package gateway

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
	userService "github.com/bsn-si/IPEHR-gateway/src/pkg/user/service"
)

type UserService interface {
	NewProcRequest(reqID, userID string, kind processing.RequestKind) (processing.RequestInterface, error)
	Register(ctx context.Context, user *model.UserCreateRequest, systemID, reqID string) (err error)
	Login(ctx context.Context, userID, systemID, password string) (err error)
	Info(ctx context.Context, userID, systemID string) (*model.UserInfo, error)
	InfoByCode(ctx context.Context, code int) (*model.UserInfo, error)
	CreateToken(userID string) (*userService.TokenDetails, error)
	ExtractToken(bearToken string) string
	VerifyAccess(userID, tokenString string) error
	VerifyToken(userID, tokenString string, tokenType userService.TokenType) (*jwt.Token, error)
	ExtractTokenMetadata(token *jwt.Token) (*userService.TokenClaims, error)
	IsTokenInBlackList(tokenRaw string) bool
	AddTokenInBlackList(tokenRaw string, expires int64)
	GetTokenHash(tokenRaw string) [32]byte
	VerifyAndGetTokenDetails(userID, accessToken, refreshToken string) (*userService.TokenDetails, error)
	GroupCreate(ctx context.Context, userID, name, description string) (string, *uuid.UUID, error)
	GroupGetByID(ctx context.Context, userID, systemID string, groupID *uuid.UUID, groupKey *chachaPoly.Key) (*model.UserGroup, error)
	GroupAddUser(ctx context.Context, userID, systemID, addUserID, addSystemID, reqID string, level access.Level, groupID *uuid.UUID) error
	GroupRemoveUser(ctx context.Context, userID, systemID, removingUserID, removeSystemID, reqID string, groupID *uuid.UUID) error
	GroupGetList(ctx context.Context, userID, systemID string) ([]*model.UserGroup, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(handlerService UserService) *UserHandler {
	return &UserHandler{
		service: handlerService,
	}
}

// Register
//
//	@Summary	Register user
//	@Description
//	@Tags		USER
//	@Accept		json
//	@Produce	json
//	@Param		EhrSystemId	header	string					false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		Request		body	model.UserCreateRequest	true	"User creation request. `role`: 0 - Patient, 1 - Doctor. Fields `Name`, `Address`, `Description`, `PictureURL` are required for Doctor role"
//	@Success	201			"Indicates that the request has succeeded and transaction about register new user has been created"
//	@Header		201			{string}	RequestID	"Request identifier"
//	@Failure	400			"The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
//	@Failure	409			"User with that userID already exist"
//	@Failure	422			"Password, systemID or role incorrect"
//	@Failure	500			"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/register [post]
func (h *UserHandler) Register(c *gin.Context) {
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

	var userCreateRequest model.UserCreateRequest
	if err = json.Unmarshal(data, &userCreateRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if ok, err := userCreateRequest.Validate(); !ok {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}

	if err := h.service.Register(c, &userCreateRequest, systemID, reqID); err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation error"})
		}

		return
	}

	c.Status(http.StatusCreated)
}

// Login
//
//	@Summary	Login user
//	@Description
//	@Tags		USER
//	@Accept		json
//	@Produce	json
//	@Param		AuthUserId	header		string					true	"UserId"
//	@Param		EhrSystemId	header		string					false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		Request		body		model.UserAuthRequest	true	"User authentication request"
//	@Success	200			{object}	model.JWT
//	@Failure	404			"User with ID not exist"
//	@Failure	422			"The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
//	@Failure	400			"Password, EhrSystemId or userID incorrect"
//	@Failure	401			"Password or userID incorrect"
//	@Failure	500			"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/login [post]
func (h *UserHandler) Login(c *gin.Context) {
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
//
//	@Summary	Logout
//	@Description
//	@Tags		USER
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header	string		true	"Bearer AccessToken"
//	@Param		AuthUserId		header	string		true	"UserId"
//	@Param		Request			body	model.JWT	true	"JWT"
//	@Success	200				"Successfully logged out"
//	@Failure	401				"User unauthorized"
//	@Failure	422				"The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
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
//
//	@Summary	Refresh JWT
//	@Description
//	@Tags		USER
//	@Accept		json
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer RefreshToken"
//	@Param		AuthUserId		header		string	true	"UserId"
//	@Success	200				{object}	model.JWT
//	@Failure	401				"User unauthorized"
//	@Failure	404				"User with ID not exist"
//	@Failure	422				"The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/refresh/ [get]
func (h *UserHandler) RefreshToken(c *gin.Context) {
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

// Info
//
//	@Summary		Get user info
//	@Description	Get information about the user by user_id
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			user_id			path		string	true	"The identifier of the requested user info"
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Param			EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200				{object}	model.UserInfo
//	@Failure		400				"`user_id` is incorrect or requested user is not a doctor"
//	@Failure		404				"User with ID not exist"
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/user/{user_id} [get]
func (h *UserHandler) Info(c *gin.Context) {
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request user_id is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	userInfo, err := h.service.Info(c, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) || errors.Is(err, errors.ErrIsNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println("service.Info error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userInfo.Role == roles.Patient.String() {
		userInfo = &model.UserInfo{
			Role:        userInfo.Role,
			TimeCreated: userInfo.TimeCreated,
			EhrID:       userInfo.EhrID,
		}
	}

	c.JSON(http.StatusOK, userInfo)
}

// Info by code
//
//	@Summary		Get doctor info by code
//	@Description	Get information about the doctor by code
//	@Tags			USER
//	@Accept			json
//	@Produce		json
//	@Param			code	path		string	true	"The pin code of the requested doctor"
//	@Success		200		{object}	model.UserInfo
//	@Failure		400		"`code` is incorrect or requested user is not a doctor"
//	@Failure		404		"User code is not exist"
//	@Failure		500		"Is returned when an unexpected error occurs while processing a request"
//	@Router			/user/code/{code} [get]
func (h *UserHandler) InfoByCode(c *gin.Context) {
	codeString := c.Param("code")

	codeInt, _ := strconv.Atoi(codeString)
	if codeInt < 0 || codeInt > 99999999 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request code must contain 8 digits"})
		return
	}

	userInfo, err := h.service.InfoByCode(c, codeInt)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println("service. InfoByCode error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userInfo.Role != roles.Doctor.String() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Info only available for users with the 'Doctor' role"})
		return
	}

	c.JSON(http.StatusOK, userInfo)
}
