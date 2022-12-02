package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

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
// @Success  201            {object} model.UserGroup "Indicates that the request has succeeded and transaction about create new user group has been created"
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

	userGroup.GroupID, err = h.service.GroupCreate(c, userID, systemID, reqID, userGroup.Name, userGroup.Description)
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

	c.JSON(http.StatusCreated, userGroup)
}

// Group get by ID
// @Summary  Get user group by ID
// @Description
// @Tags     USER
// @Produce  json
// @Param    Authorization  header  string           true  "Bearer AccessToken"
// @Param    AuthUserId     header  string           true  "UserId"
// @Param    EhrSystemId    header  string           true  "The identifier of the system, typically a reverse domain identifier"
// @Success  200            {object}  model.UserGroup
// @Failure  400            "The request could not be understood by the server due to incorrect syntax."
// @Failure  403            "Is returned when userID does not have access to requested group"
// @Failure  404            "Is returned when groupID does not exist"
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /user/group/{group_id} [get]
func (h *UserHandler) GroupGetByID(c *gin.Context) {
	gID := c.Param("group_id")
	if gID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id is empty"})
		return
	}

	groupID, err := uuid.Parse(gID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group_id must be UUID"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	userGroup, err := h.service.GroupGetByID(c, userID, &groupID)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println("GroupCreate error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userGroup)
}
