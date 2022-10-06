package api

import (
	"encoding/json"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/user"
	"hms/gateway/pkg/errors"
	"io"
	"log"
	"net/http"

	"hms/gateway/pkg/docs/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(docService *service.DefaultDocumentService) *UserHandler {
	return &UserHandler{
		service: user.NewUserService(docService),
	}
}

// Register
// @Summary      Register user
// @Description
// @Tags     REQUEST
// @Accept   json
// @Produce  json
// @Param    userID     body    string    true  "UserId UUID"
// @Param    systemID   body    string    true  "The identifier of the system, typically a reverse domain identifier"
// @Param    password   body    string	  true  "Password"
// @Param    role       body    string    true  ""
// TODO can users register by themselves, or does it have to be an already authorized user?
// @Success  201         "Indicates that the request has succeeded and a new user has been created"
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

	var request model.UserCreateRequest

	if err = json.Unmarshal(data, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Request validation error"})
		return
	}

	err = h.service.Register(&request)
	if err != nil {
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
