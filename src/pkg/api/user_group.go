package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/user/model"
)

// Group create
// @Summary  User group create
// @Description
// @Tags     USER
// @Accept   json
// @Param    Authorization  header  string           true  "Bearer AccessToken"
// @Param    AuthUserId     header  string           true  "UserId UUID"
// @Param    EhrSystemId    header  string           true  "The identifier of the system, typically a reverse domain identifier"
// @Param    Request        body    model.UserGroup  true  "User group"
// @Success  201            "Indicates that the request has succeeded and transaction about create new user group has been created"
// @Header   201            {string}  RequestID  "Request identifier"
// @Failure  400            "The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
// @Failure  404            "User with ID not exist"
// @Failure  409            "Group with that Name already exist"
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/group [post]
func (h *UserHandler) GroupCreate(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header EhrSystemId is empty"})
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}

	var userGroup model.UserGroup

	err = json.Unmarshal(data, &userGroup)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if userGroup.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request group name is empty"})
		return
	}

	reqID := c.GetString("reqID")

	err = h.service.GroupCreate(c, userID, systemID, reqID, userGroup.Name, userGroup.Description)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		log.Println("GroupCreate error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusCreated)
}
