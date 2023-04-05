package gateway

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

// Group create
//
//	@Summary	User group create
//	@Description
//	@Tags		USER_GROUP
//	@Accept		json
//	@Param		Authorization	header		string			true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string			true	"UserId"
//	@Param		Request			body		model.UserGroup	true	"User group"
//	@Success	201				{object}	model.UserGroup	"Indicates that the request has succeeded and transaction about create new user group has been created"
//	@Header		201				{string}	RequestID		"Request identifier"
//	@Failure	400				"The request could not be understood by the server due to incorrect syntax. The client SHOULD NOT repeat the request without modifications."
//	@Failure	404				"User with ID not exist"
//	@Failure	409				"Group with that Name already exist"
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/group [post]
func (h *UserHandler) GroupCreate(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
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

	var txHash string

	txHash, userGroup.GroupID, err = h.service.GroupCreate(c, userID, userGroup.Name, userGroup.Description)
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

	reqID := c.GetString("reqID")

	procRequest, err := h.service.NewProcRequest(reqID, userID, processing.RequestUserGroupCreate)
	if err != nil {
		log.Println("Proc.NewRequest error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	procRequest.AddEthereumTx(processing.TxUserGroupCreate, txHash)

	if err := procRequest.Commit(); err != nil {
		log.Println("UserGroup create procRequest.Commit error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, userGroup)
}

// Group get by ID
//
//	@Summary	Get user group by ID
//	@Description
//	@Tags		USER_GROUP
//	@Produce	json
//	@Param		group_id		path		string	true	"User group identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		Authorization	header		string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string	true	"UserId"
//	@Param		EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success	200				{object}	model.UserGroup
//	@Failure	400				"The request could not be understood by the server due to incorrect syntax."
//	@Failure	403				"Is returned when userID does not have access to requested group"
//	@Failure	404				"Is returned when groupID does not exist"
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/group/{group_id} [get]
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

	systemID := c.GetString("ehrSystemID")

	userGroup, err := h.service.GroupGetByID(c, userID, systemID, &groupID, nil)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println("GroupGetByID error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, userGroup)
}

// Group add user
//
//	@Summary	Adding a user to a group
//	@Description
//	@Tags		USER_GROUP
//	@Accept		json
//	@Param		Authorization	header	string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header	string	true	"UserId"
//	@Param		EhrSystemId		header	string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		group_id		path	string	true	"The identifier of the user group"
//	@Param		user_id			path	string	true	"The identifier of the user to be added"
//	@Param		access_level	path	string	true	"Access Level. One of `admin` or `read`"
//	@Success	200				""
//	@Header		200				{string}	RequestID	"Request identifier"
//	@Failure	400				"The request could not be understood by the server due to incorrect syntax."
//	@Failure	403				"Authentication required or user does not have access to change the group"
//	@Failure	404				"Group or adding user is not exist"
//	@Failure	409				"The user is already a member of a group"
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/group/{group_id}/user_add/{user_id}/{access_level} [put]
func (h *UserHandler) GroupAddUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

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

	addUserID := c.Param("user_id")
	if addUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is empty"})
		return
	}

	level := access.LevelFromString(c.Param("access_level"))
	if level == access.Unknown {
		c.JSON(http.StatusBadRequest, gin.H{"error": "access_level is incorrect"})
		return
	}

	reqID := c.GetString("reqID")

	//TODO split systemID and adding user SystemID
	err = h.service.GroupAddUser(c, userID, systemID, addUserID, systemID, reqID, level, &groupID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			c.AbortWithStatus(http.StatusConflict)
			return
		} else if errors.Is(err, errors.ErrAccessDenied) {
			log.Println("GroupAddUser error: ", err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		log.Println("GroupAddUser error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// Group remove user
//
//	@Summary	Removing a user from a group
//	@Description
//	@Tags		USER_GROUP
//	@Accept		json
//	@Param		Authorization	header	string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header	string	true	"UserId"
//	@Param		EhrSystemId		header	string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		group_id		path	string	true	"The identifier of the user group"
//	@Param		user_id			path	string	true	"The identifier of the user to be removed"
//	@Success	200				""
//	@Header		200				{string}	RequestID	"Request identifier"
//	@Failure	400				"The request could not be understood by the server due to incorrect syntax."
//	@Failure	403				"Authentication required or user does not have access to change the group"
//	@Failure	404				"Group or adding user is not exist or `user_id` is not the member of the group"
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/group/{group_id}/user_remove/{user_id} [post]
func (h *UserHandler) GroupRemoveUser(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

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

	removeUserID := c.Param("user_id")
	if removeUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is empty"})
		return
	}

	reqID := c.GetString("reqID")

	err = h.service.GroupRemoveUser(c, userID, systemID, removeUserID, systemID, reqID, &groupID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrAccessDenied) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		log.Println("GroupRemoveUser error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}

// Group get list
//
//	@Summary	Get a list of user groups
//	@Description
//	@Tags		USER_GROUP
//	@Produce	json
//	@Param		Authorization	header		string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string	true	"UserId"
//	@Param		EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success	200				{object}	[]model.UserGroup
//	@Failure	400				"The request could not be understood by the server due to incorrect syntax."
//	@Failure	404				"Is returned when groupID does not exist"
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/user/group [get]
func (h *UserHandler) GroupGetList(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	groupList, err := h.service.GroupGetList(c, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println("GroupGetList error: ", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(groupList) == 0 {
		groupList = make([]*model.UserGroup, 0)
	}

	c.JSON(http.StatusOK, groupList)
}
