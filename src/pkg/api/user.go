package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/user/service"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(cfg *config.Config, infra *infrastructure.Infra, p *proc.Proc) *UserHandler {
	return &UserHandler{
		service: service.NewUserService(cfg, infra, p),
	}
}

// Register
// @Summary  Register user
// @Description
// @Tags     USER
// @Accept   json
// @Produce  json
// @Param    EhrSystemId  header  string                   true  "The identifier of the system, typically a reverse domain identifier"
// @Param    Request      body    model.UserCreateRequest  true  "User creation request"
// @Success  201          "Indicates that the request has succeeded and transaction about register new user has been created"
// @Header   201          {string}  RequestID  "Request identifier"
// @Failure  400          "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  409          "User with that userID already exist"
// @Failure  422          "Password, systemID or role incorrect"
// @Failure  500          "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/register [post]
func (h UserHandler) Register(c *gin.Context) {
	// TODO can users register by themselves, or does it have to be an already authorized user?
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	reqID := c.GetString("reqID")
	if reqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requestId is empty"})
		return
	}

	// Label 'register' to allow the status of a user registration request without auth
	// Used in auth exceptions middleware
	reqID += "_register"
	c.Header("RequestId", reqID)

	systemID := c.GetString("ehrSystemID")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EhrSystemId required"})
		return
	}

	var userCreateRequest model.UserCreateRequest
	if err = json.Unmarshal(data, &userCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if ok, err := userCreateRequest.Validate(); !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	procRequest, err := h.service.Proc.NewRequest(reqID, userCreateRequest.UserID, "", proc.RequestUserRegister)
	if err != nil {
		log.Println("User register NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = h.service.Register(c, procRequest, &userCreateRequest, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation error"})
		}

		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("User register procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}

// Login
// @Summary  Login user
// @Description
// @Tags     USER
// @Accept   json
// @Produce  json
// @Param    AuthUserId   header    string                 true  "UserId UUID"
// @Param    EhrSystemId  header    string                 true  "The identifier of the system, typically a reverse domain identifier"
// @Param    Request      body      model.UserAuthRequest  true  "User authentication request"
// @Success  200          {object}  model.JWT
// @Failure  404          "User with ID not exist"
// @Failure  422          "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  400          "Password, EhrSystemId or userID incorrect"
// @Failure  401          "Password or userID incorrect"
// @Failure  500          "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/login [post]
func (h UserHandler) Login(c *gin.Context) {
	// TODO add timeout between attempts, we dont need password brute force
	var u model.UserAuthRequest
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	if ok, err := u.Validate(); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	tokenString := c.Request.Header.Get("Authorization")
	userID := c.Request.Header.Get("AuthUserId")
	systemID := c.GetString("ehrSystemID")

	if tokenString != "" && userID != "" {
		c.JSON(http.StatusUnprocessableEntity, "You are already authorised, please use logout")
		return
	}

	err := h.service.Login(c, u.UserID, systemID, u.Password)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		log.Println(err)
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ts, err := h.service.CreateToken(u.UserID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, &model.JWT{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	})
}

// Logout
// @Summary  Logout
// @Description
// @Tags     USER
// @Accept   json
// @Produce  json
// @Param    Authorization  header  string     true  "Bearer AccessToken"
// @Param    AuthUserId     header  string     true  "UserId - UUID"
// @Param    Request        body    model.JWT  true  "JWT"
// @Success  200            "Successfully logged out"
// @Failure  401            "User unauthorized"
// @Failure  422            "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/logout/ [post]
func (h UserHandler) Logout(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	userID := c.Request.Header.Get("AuthUserId")

	tokenString = h.service.ExtractToken(tokenString)

	if tokenString == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var jwt model.JWT
	if err := c.ShouldBindJSON(&jwt); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	td, err := h.service.VerifyAndGetTokenDetails(userID, tokenString, jwt.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	h.service.AddTokenInBlackList(td.AccessToken, td.AtExpires)
	h.service.AddTokenInBlackList(td.RefreshToken, td.RtExpires)

	c.JSON(http.StatusOK, "Successfully logged out")
}

// RefreshToken
// @Summary  Refresh JWT
// @Description
// @Tags     USER
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string  true  "Bearer RefreshToken"
// @Param    AuthUserId     header    string  true  "UserId - UUID"
// @Param    EhrSystemId    header    string  true  "The identifier of the system, typically a reverse domain identifier"
// @Success  200            {object}  model.JWT
// @Failure  401            "User unauthorized"
// @Failure  404            "User with ID not exist"
// @Failure  422            "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/refresh/ [get]
func (h UserHandler) RefreshToken(c *gin.Context) {
	userID := c.Request.Header.Get("AuthUserId")
	tokenString := c.Request.Header.Get("Authorization")

	tokenString = h.service.ExtractToken(tokenString)

	if userID == "" || tokenString == "" {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	td, err := h.service.VerifyAndGetTokenDetails(userID, "", tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	h.service.AddTokenInBlackList(td.RefreshToken, td.RtExpires)

	ts, err := h.service.CreateToken(userID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	c.JSON(http.StatusOK, &model.JWT{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	})
}
