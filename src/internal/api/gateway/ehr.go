package gateway

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	docGroupService "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/docGroup"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/ehr"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/groupAccess"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	userService "github.com/bsn-si/IPEHR-gateway/src/pkg/user/service"
)

type EhrHandler struct {
	service *ehr.Service
	baseURL string
}

func NewEhrHandler(docSvc *service.DefaultDocumentService, userSvc *userService.Service, docGroupSvc *docGroupService.Service, gaSvc *groupAccess.Service, baseURL string) *EhrHandler {
	return &EhrHandler{
		service: ehr.NewService(docSvc, userSvc, docGroupSvc, gaSvc),
		baseURL: baseURL,
	}
}

// Create
//
//	@Summary		Create EHR
//	@Description	Create a new EHR with an auto-generated identifier.
//	@Description	An EHR_STATUS resource needs to be always created and committed in the new EHR. This resource MAY be also supplied by the client as the request body. If not supplied, a default EHR_STATUS will be used by the service with following attributes:
//	@Description	- `is_queryable`: true
//	@Description	- `is_modifiable`: true
//	@Description	- `subject`: a PARTY_SELF object
//	@Description
//	@Description	All other required EHR attributes and resources will be automatically created as needed by the [EHR creation semantics](https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_creation_semantics).
//	@Tags			EHR
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string					true	"UserId"
//	@Param			EhrSystemId		header		string					false	"The identifier of the system, typically a reverse domain identifier"
//	@Param			GroupAccessId	header		string					false	"GroupAccessId - UUID. If not specified, the default access group will be used."
//	@Param			Prefer			header		string					true	"The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
//	@Param			Request			body		model.EhrCreateRequest	true	"Query Request"
//	@Success		201				{object}	model.EhrSummary
//	@Header			201				{string}	Location	"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Header			201				{string}	ETag		"ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Header			201				{string}	RequestID	"Request identifier"
//	@Failure		400				"Is returned when the request body (if provided)  could      not  be  parsed."
//	@Failure		409				"Unable to create a new EHR due to a conflict with an already existing EHR with the same subject id, namespace pair."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr [post]
func (h *EhrHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var request model.EhrCreateRequest

	if err = json.Unmarshal(data, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	// Checking EHR does not exist?
	// проверить чему будет равен ehrUUID
	ehrUUID, err := h.service.Infra.Index.GetEhrUUIDByUserID(c, userID, systemID)
	switch {
	case err == nil && ehrUUID != nil:
		c.AbortWithStatus(http.StatusConflict)
		return
	case err != nil && !errors.Is(err, errors.ErrIsNotExist):
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	ehrUUIDnew := uuid.New()
	reqID := c.GetString("reqID")

	groupAccessUUID := h.service.GroupAccess.Default()

	if c.GetHeader("GroupAccessId") != "" {
		UUID, err := uuid.Parse(c.GetHeader("GroupAccessId"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "GroupAccessId parsing error"})
			return
		}

		groupAccessUUID = &UUID
	}

	procRequest, err := h.service.Doc.Proc.NewRequest(reqID, userID, ehrUUIDnew.String(), processing.RequestEhrCreate)
	if err != nil {
		log.Println("EHR create NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// EHR document creating
	doc, err := h.service.EhrCreate(c, userID, systemID, &ehrUUIDnew, groupAccessUUID, &request, procRequest)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}

		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("EHR create procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.respondWithDocOrHeaders(doc, c)
}

// CreateWithID
//
//	@Summary		Create EHR with id
//	@Description	Create a new EHR with the specified ehr_id identifier.
//	@Description	The value of the ehr_id unique identifier MUST be valid HIER_OBJECT_ID value. It is strongly RECOMMENDED that an UUID always be used for this.
//	@Description	An EHR_STATUS resource needs to be always created and committed in the new EHR. This resource MAY be also supplied by the client as the request body. If not supplied, a default EHR_STATUS will be used by the service with following attributes:
//	@Description	- `is_queryable`: true
//	@Description	- `is_modifiable`: true
//	@Description	- `subject`: a PARTY_SELF object
//	@Description
//	@Description	All other required EHR attributes and resources will be automatically created as needed by the [EHR creation semantics](https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_creation_semantics).
//	@Tags			EHR
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string					true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string					true	"UserId"
//	@Param			EhrSystemId		header		string					false	"The identifier of the system, typically a reverse domain identifier"
//	@Param			Prefer			header		string					true	"The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
//	@Param			ehr_id			path		string					true	"An UUID as a user specified EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param			Request			body		model.EhrCreateRequest	true	"Query Request"
//	@Success		201				{object}	model.EhrSummary
//	@Header			201				{string}	Location	"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Header			201				{string}	ETag		"ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Header			201				{string}	RequestID	"Request identifier"
//	@Failure		400				"Is returned when the request body (if provided)  could      not  be  parsed."
//	@Failure		409				"Unable to create a new EHR due to a conflict with an already existing EHR. Can happen when the supplied ehr_id is already used by an existing EHR."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr/{ehr_id} [put]
func (h *EhrHandler) CreateWithID(c *gin.Context) {
	ehrID := c.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ehr_id"})
		return
	}

	request := model.EhrCreateRequest{}
	if err = json.NewDecoder(c.Request.Body).Decode(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body parsing error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	reqID := c.GetString("reqID")
	systemID := c.GetString("ehrSystemID")

	// Checking EHR does not exist
	_ehrUUID, err := h.service.Infra.Index.GetEhrUUIDByUserID(c, userID, systemID)
	switch {
	case err == nil && _ehrUUID != nil:
		c.AbortWithStatus(http.StatusConflict)
		return
	case err != nil && !errors.Is(err, errors.ErrIsNotExist):
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	groupAccessUUID := h.service.GroupAccess.Default()

	if c.GetHeader("GroupAccessId") != "" {
		UUID, err := uuid.Parse(c.GetHeader("GroupAccessId"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "GroupAccessId parsing error"})
			return
		}

		groupAccessUUID = &UUID
	}

	procRequest, err := h.service.Doc.Proc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestEhrCreateWithID)
	if err != nil {
		log.Println("Ehr createWithID NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	// EHR document creating
	newDoc, err := h.service.EhrCreateWithID(c, userID, systemID, &ehrUUID, groupAccessUUID, &request, procRequest)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		}

		log.Println("EhrCreateWithID error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})

		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("EHR createWithID procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.respondWithDocOrHeaders(newDoc, c)
}

// GetByID
//
//	@Summary		Get EHR summary by id
//	@Description	Retrieve the EHR with the specified ehr_id
//	@Tags			EHR
//	@Accept			json
//	@Produce		json
//	@Param			ehr_id			path		string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Param			EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200				{object}	model.EhrSummary
//	@Success		202				"Is returned when the request is still being processed"
//	@Failure		400				"Is returned when userID is empty"
//	@Failure		404				"Is returned when an EHR with ehr_id does not exist."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr/{ehr_id} [get]
func (h *EhrHandler) GetByID(c *gin.Context) {
	ehrID := c.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ehr_id"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	docDecrypted, err := h.service.GetByID(c, userID, systemID, &ehrUUID)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		c.AbortWithStatus(http.StatusAccepted)
		return
	} else if err != nil {
		log.Printf("service.GetByID error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document getting from storage error"})
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

// GetBySubjectIDAndNamespace
//
//	@Summary		Get EHR summary by subject id
//	@Description	Retrieve the EHR with the specified subject_id and subject_namespace.
//	@Description	These subject parameters will be matched against EHR’s
//	@Description	EHR_STATUS.subject.external_ref.id.value and EHR_STATUS.subject.external_ref.namespace values.
//	@Tags			EHR
//	@Accept			json
//	@Produce		json
//	@Param			subject_id			query		string	true	"subject id. Example: ins01"
//	@Param			subject_namespace	query		string	true	"id namespace. Example: examples"
//	@Param			Authorization		header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId			header		string	true	"UserId"
//	@Param			EhrSystemId			header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200					{object}	model.EhrSummary
//	@Success		202					"Is returned when the request is still being processed"
//	@Failure		400					"Is returned when userID is empty"
//	@Failure		404					"Is returned when an EHR with ehr_id does not exist."
//	@Router			/ehr [get]
func (h *EhrHandler) GetBySubjectIDAndNamespace(c *gin.Context) {
	subjectID := c.Query("subject_id")
	namespace := c.Query("subject_namespace")

	if subjectID == "" || namespace == "" {
		log.Println("Subject data is not filled correctly")
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	docDecrypted, err := h.service.GetDocBySubject(c, userID, systemID, subjectID, namespace)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		c.AbortWithStatus(http.StatusAccepted)
		return
	} else if err != nil {
		log.Println("Can't get document by subject", err)
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

func (h *EhrHandler) respondWithDocOrHeaders(doc *model.EHR, c *gin.Context) {
	c.Header("Location", h.baseURL+"/v1/ehr/"+doc.EhrID.Value)
	c.Header("ETag", doc.EhrID.Value)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
