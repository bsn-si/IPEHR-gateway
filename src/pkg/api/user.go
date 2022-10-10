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

func (h UserHandler) Login(c *gin.Context) {
	// TODO if we have AuthUserId or Barrier should be say 'gerrarahea!)???'
}

func (h UserHandler) Refresh(c *gin.Context) {
	// TODO refresh token exp
}

func (h UserHandler) LoginOut(c *gin.Context) {
	// TODO exit
}
