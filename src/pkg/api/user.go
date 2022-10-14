package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/infrastructure"
	"hms/gateway/pkg/user/service"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(cfg *config.Config, infra *infrastructure.Infra) *UserHandler {
	return &UserHandler{
		service: service.NewUserService(cfg, infra),
	}
}

// Register
// @Summary      Register user
// @Description
// @Tags     REQUEST
// @Accept   json
// @Produce  json
// @Param    userID       body    string  true  "UserId UUID"
// @Param    password     body    string  true  "Password"
// @Param    role         body    string  true  ""
// @Param    EhrSystemId  header  string  true  "The identifier of the system, typically a reverse domain identifier"
// TODO can users register by themselves, or does it have to be an already authorized user?
// @Success  201         "Indicates that the request has succeeded and transaction about register new user has been created"
// @Header   201 {string}  RequestID  "Request identifier"
// @Failure  400         "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  409         "User with that userId already exist"
// @Failure  422         "Password, systemID or role incorrect"
// @Failure  500         "Is returned when an unexpected error occurs while processing a request"
// @Router   /requests/ [get]
func (h UserHandler) Register(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	reqID := c.MustGet("reqId").(string)
	_ = c.MustGet("ehrSystemID").(base.EhrSystemID)

	if reqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requestId is empty"})
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

	err = h.service.Register(c, procRequest, &userCreateRequest)
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

// TODO fill comments
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

	if tokenString != "" && userID != "" {
		c.JSON(http.StatusUnprocessableEntity, "You are already authorised, please use logout")
		return
	}

	err := h.service.Login(c, &u)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}

		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	ts, err := h.service.CreateToken(u.UserID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := h.service.CreateAuth(u.UserID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	c.JSON(http.StatusCreated, &model.JWT{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	})
}

// TODO fill comments
func (h UserHandler) Logout(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	userID := c.Request.Header.Get("AuthUserId")

	tokenString = h.service.ExtractToken(tokenString)

	if tokenString == "" || userID == "" {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	token, err := h.service.VerifyToken(userID, tokenString, false)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	metadata, err := h.service.ExtractTokenMetadata(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	delErr := h.service.DeleteToken(metadata.Uuid)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

// TODO fill comments
func (h UserHandler) RefreshToken(c *gin.Context) {
	userID := c.Request.Header.Get("AuthUserId")
	tokenString := c.Request.Header.Get("Authorization")

	tokenString = h.service.ExtractToken(tokenString)

	if userID == "" || tokenString == "" {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	var jwt model.JWT
	if err := c.ShouldBindJSON(&jwt); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	token, err := h.service.VerifyToken(userID, jwt.RefreshToken, true)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}

	metadata, err := h.service.ExtractTokenMetadata(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	delErr := h.service.DeleteToken(metadata.Uuid)
	if delErr != nil {
		c.JSON(http.StatusUnauthorized, delErr.Error())
		return
	}

	ts, err := h.service.CreateToken(metadata.UserId)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	saveErr := h.service.CreateAuth(metadata.UserId, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}

	c.JSON(http.StatusCreated, &model.JWT{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	})
}
