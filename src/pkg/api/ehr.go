package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/service/processing"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type EhrHandler struct {
	service *ehr.Service
	baseURL string
}

func NewEhrHandler(docService *service.DefaultDocumentService, baseURL string) *EhrHandler {
	return &EhrHandler{
		service: ehr.NewService(docService),
		baseURL: baseURL,
	}
}

// Create
// @Summary      Create EHR
// @Description  Create a new EHR with an auto-generated identifier.
// @Description  An EHR_STATUS resource needs to be always created and committed in the new EHR. This resource MAY be also supplied by the client as the request body. If not supplied, a default EHR_STATUS will be used by the service with following attributes:
// @Description  - `is_queryable`: true
// @Description  - `is_modifiable`: true
// @Description  - `subject`: a PARTY_SELF object
// @Description
// @Description  All other required EHR attributes and resources will be automatically created as needed by the [EHR creation semantics](https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_creation_semantics).
// @Tags         EHR
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer AccessToken"
// @Param        AuthUserId     header    string                  true  "UserId UUID"
// @Param        EhrSystemId    header    string                  true  "The identifier of the system, typically a reverse domain identifier"
// @Param        Prefer         header    string                  true  "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        Request        body      model.EhrCreateRequest  true  "Query Request"
// @Success      201            {object}  model.EhrSummary
// @Header       201            {string}  Location   "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201            {string}  ETag       "ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201            {string}  RequestID  "Request identifier"
// @Failure      400            "Is returned when the request body (if provided) could not be parsed."
// @Failure      409            "Unable to create a new EHR due to a conflict with an already existing EHR with the same subject id, namespace pair."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr [post]
func (h *EhrHandler) Create(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

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
	ehrUUID, err := h.service.Infra.Index.GetEhrUUIDByUserID(c, userID)
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
	ehrSystemID := c.GetString("ehrSystemID")
	reqID := c.GetString("reqID")

	procRequest, err := h.service.Proc.NewRequest(reqID, userID, ehrUUIDnew.String(), processing.RequestEhrCreate)
	if err != nil {
		log.Println("EHR create NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// EHR document creating
	doc, err := h.service.EhrCreate(c, userID, &ehrUUIDnew, ehrSystemID, &request, procRequest)
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
// @Summary      Create EHR with id
// @Description  Create a new EHR with the specified ehr_id identifier.
// @Description  The value of the ehr_id unique identifier MUST be valid HIER_OBJECT_ID value. It is strongly RECOMMENDED that an UUID always be used for this.
// @Description  An EHR_STATUS resource needs to be always created and committed in the new EHR. This resource MAY be also supplied by the client as the request body. If not supplied, a default EHR_STATUS will be used by the service with following attributes:
// @Description  - `is_queryable`: true
// @Description  - `is_modifiable`: true
// @Description  - `subject`: a PARTY_SELF object
// @Description
// @Description  All other required EHR attributes and resources will be automatically created as needed by the [EHR creation semantics](https://specifications.openehr.org/releases/RM/latest/ehr.html#_ehr_creation_semantics).
// @Tags         EHR
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                  true  "Bearer AccessToken"
// @Param        AuthUserId     header    string                  true  "UserId UUID"
// @Param        EhrSystemId    header    string                  true  "The identifier of the system, typically a reverse domain identifier"
// @Param        Prefer         header    string                  true  "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        ehr_id         path      string                  true  "An UUID as a user specified EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        Request        body      model.EhrCreateRequest  true  "Query Request"
// @Success      201            {object}  model.EhrSummary
// @Header       201            {string}  Location   "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201            {string}  ETag       "ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201            {string}  RequestID  "Request identifier"
// @Failure      400            "Is returned when the request body (if provided) could not be parsed."
// @Failure      409            "Unable to create a new EHR due to a conflict with an already existing EHR. Can happen when the supplied ehr_id is already used by an existing EHR."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id} [put]
func (h *EhrHandler) CreateWithID(c *gin.Context) {
	ehrID := c.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect ehr_id"})
		return
	}

	ehrSystemID := c.GetString("ehrSystemID")

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

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	// Checking EHR does not exist
	_ehrUUID, err := h.service.Infra.Index.GetEhrUUIDByUserID(c, userID)
	switch {
	case err == nil && _ehrUUID != nil:
		c.AbortWithStatus(http.StatusConflict)
		return
	case err != nil && !errors.Is(err, errors.ErrIsNotExist):
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	reqID := c.GetString("reqID")

	procRequest, err := h.service.Proc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestEhrCreateWithID)
	if err != nil {
		log.Println("Ehr createWithID NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	// EHR document creating
	newDoc, err := h.service.EhrCreateWithID(c, userID, &ehrUUID, ehrSystemID, &request, procRequest)
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
// @Summary      Get EHR summary by id
// @Description  Retrieve the EHR with the specified ehr_id
// @Tags         EHR
// @Accept       json
// @Produce      json
// @Param        ehr_id         path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        Authorization  header    string  true  "Bearer AccessToken"
// @Param        AuthUserId     header    string  true  "UserId UUID"
// @Param        EhrSystemId    header    string  true  "The identifier of the system, typically a reverse domain identifier"
// @Success      200            {object}  model.EhrSummary
// @Success      202            "Is returned when the request is still being processed"
// @Failure      400            "Is returned when userID is empty"
// @Failure      404            "Is returned when an EHR with ehr_id does not exist."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id} [get]
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

	// Getting docStorageID
	docMeta, err := h.service.Infra.Index.GetDocLastByType(c, &ehrUUID, types.Ehr)
	if err != nil {
		log.Println("GetLastDocIndexByType", "ehrId", ehrID, err)
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	// Getting doc from storage
	docDecrypted, err := h.service.GetDocFromStorageByID(c, userID, docMeta.Cid(), ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		c.AbortWithStatus(http.StatusAccepted)
		return
	} else if err != nil {
		log.Printf("GetDocFromStorageByID error: %v\ndocMeta: %+v", err, docMeta)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document getting from storage error"})
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

// GetBySubjectIDAndNamespace
// @Summary      Get EHR summary by subject id
// @Description  Retrieve the EHR with the specified subject_id and subject_namespace.
// @Description  These subject parameters will be matched against EHR’s
// @Description  EHR_STATUS.subject.external_ref.id.value and EHR_STATUS.subject.external_ref.namespace values.
// @Tags         EHR
// @Accept       json
// @Produce      json
// @Param        subject_id         query     string  true  "subject id. Example: ins01"
// @Param        subject_namespace  query     string  true  "id namespace. Example: examples"
// @Param        Authorization      header    string  true  "Bearer AccessToken"
// @Param        AuthUserId         header    string  true  "UserId UUID"
// @Param        EhrSystemId        header    string  true  "The identifier of the system, typically a reverse domain identifier"
// @Success      200                {object}  model.EhrSummary
// @Success      202                "Is returned when the request is still being processed"
// @Failure      400                "Is returned when userID is empty"
// @Failure      404                "Is returned when an EHR with ehr_id does not exist."
// @Router       /ehr [get]
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

	docDecrypted, err := h.service.GetDocBySubject(c, userID, subjectID, namespace)
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
