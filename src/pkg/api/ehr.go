package api

import (
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type EhrHandler struct {
	cfg     *config.Config
	service *ehr.Service
}

func NewEhrHandler(docService *service.DefaultDocumentService, cfg *config.Config) *EhrHandler {
	return &EhrHandler{
		cfg:     cfg,
		service: ehr.NewService(docService),
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
// @Param        AuthUserId  header    string                  true  "UserId UUID"
// @Param        EhrSystemId header    string                  true  "The identifier of the system, typically a reverse domain identifier"
// @Param        Prefer      header    string                  true  "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        Request     body      model.EhrCreateRequest  true  "Query Request"
// @Success      201         {object}  model.EhrSummary
// @Header       201         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201         {string}  ETag      "ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Failure      400         "Is returned when the request body (if provided) could not be parsed."
// @Failure      409         "Unable to create a new EHR due to a conflict with an already existing EHR with the same subject id, namespace pair."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr [post]
func (h EhrHandler) Create(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
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

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if !errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	ehrSystemID := c.MustGet("ehrSystemID").(base.EhrSystemID)

	// EHR document creating
	doc, err := h.service.EhrCreate(userID, ehrSystemID, &request)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}

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
// @Param        AuthUserId  header    string                  true  "UserId UUID"
// @Param        EhrSystemId header    string                  true  "The identifier of the system, typically a reverse domain identifier"
// @Param        Prefer      header    string                  true  "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        ehr_id      path      string                  true  "An UUID as a user specified EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        Request     body      model.EhrCreateRequest  true  "Query Request"
// @Success      201         {object}  model.EhrSummary
// @Header       201         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398"
// @Header       201         {string}  ETag      "ehr_id of created document. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Failure      400         "Is returned when the request body (if provided) could not be parsed."
// @Failure      409         "Unable to create a new EHR due to a conflict with an already existing EHR. Can happen when the supplied ehr_id is already used by an existing EHR."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id} [put]
func (h EhrHandler) CreateWithID(c *gin.Context) {
	ehrID := c.Param("ehrid")

	ehrSystemID := c.MustGet("ehrSystemID").(base.EhrSystemID)

	if !h.service.Doc.ValidateID(ehrID, ehrSystemID, types.Ehr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}

	// Checking EHR does not exist
	doc, err := h.service.Doc.DocsIndex.GetLastByType(ehrID, types.Ehr)
	if err != nil && !errors.Is(err, errors.ErrIsNotExist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EHR retrieve error"})
		return
	}

	if doc != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
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

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if !errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// EHR document creating
	newDoc, err := h.service.EhrCreateWithID(userID, ehrID, ehrSystemID, &request)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})

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
// @Param        ehr_id      path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        AuthUserId         header    string  true  "UserId UUID"
// @Param        EhrSystemId        header    string  true  "The identifier of the system, typically a reverse domain identifier"
// @Success      200                {object}  model.EhrSummary
// @Failure      400                "Is returned when userId is empty"
// @Failure      404                "Is returned when an EHR with ehr_id does not exist."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id} [get]
func (h EhrHandler) GetByID(c *gin.Context) {
	ehrID := c.Param("ehrid")
	ehrSystemID := c.MustGet("ehrSystemID").(base.EhrSystemID)

	if !h.service.Doc.ValidateID(ehrID, ehrSystemID, types.Ehr) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Getting docStorageID
	doc, err := h.service.Doc.DocsIndex.GetLastByType(ehrID, types.Ehr)
	if err != nil {
		log.Println("GetLastDocIndexByType", "ehrId", ehrID, err)
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	// Getting doc from storage
	docDecrypted, err := h.service.Doc.GetDocFromStorageByID(userID, doc.StorageID, []byte(ehrID))
	if err != nil {
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
// @Param        AuthUserId  header    string  true  "UserId UUID"
// @Success      200         {object}  model.EhrSummary
// @Failure      400         "Is returned when userId is empty"
// @Failure      404         "Is returned when an EHR with ehr_id does not exist."
// @Router       /ehr [get]
func (h EhrHandler) GetBySubjectIDAndNamespace(c *gin.Context) {
	subjectID := c.Query("subject_id")
	namespace := c.Query("subject_namespace")

	if subjectID == "" || namespace == "" {
		log.Println("Subject data is not filled correctly")
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docDecrypted, err := h.service.GetDocBySubject(userID, subjectID, namespace)
	if err != nil {
		log.Println("Can't get document by subject", err)
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

func (h *EhrHandler) respondWithDocOrHeaders(doc *model.EHR, c *gin.Context) {
	c.Header("Location", h.cfg.BaseURL+"/v1/ehr/"+doc.EhrID.Value)
	c.Header("ETag", doc.EhrID.Value)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
