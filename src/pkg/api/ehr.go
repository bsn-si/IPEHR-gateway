package api

import (
	"encoding/json"
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
	*ehr.EhrService
}

func NewEhrHandler(docService *service.DefaultDocumentService, cfg *config.Config) *EhrHandler {
	return &EhrHandler{
		ehr.NewEhrService(docService, cfg),
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

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.Doc.EhrsIndex.Get(userId)
	if !errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// EHR document creating
	doc, err := h.EhrCreate(userId, &request)
	if err != nil {
		if errors.Is(err, errors.AlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}
		return
	}

	h.respondWithDocOrHeaders(doc, c)
}

// CreateWithId
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
func (h EhrHandler) CreateWithId(c *gin.Context) {
	ehrId := c.Param("ehrid")

	// Checking EHR does not exist
	doc, err := h.Doc.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil && !errors.Is(err, errors.IsNotExist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EHR retrieve error"})
		return
	}
	if doc != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		return
	}

	if h.Doc.ValidateId(ehrId, types.EHR) == false {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
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

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.Doc.EhrsIndex.Get(userId)
	if !errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// EHR document creating
	newDoc, err := h.EhrCreateWithId(userId, ehrId, &request)
	if err != nil {
		if errors.Is(err, errors.AlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}
		return
	}

	h.respondWithDocOrHeaders(newDoc, c)
}

// GetById
// @Summary      Get EHR summary by id
// @Description  Retrieve the EHR with the specified ehr_id
// @Tags         EHR
// @Accept       json
// @Produce      json
// @Param        ehr_id      path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        AuthUserId         header    string  true  "UserId UUID"
// @Success      200                {object}  model.EhrSummary
// @Failure      400                "Is returned when userId is empty"
// @Failure      404                "Is returned when an EHR with ehr_id does not exist."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id} [get]
func (h EhrHandler) GetById(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Getting docStorageId
	doc, err := h.Doc.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil {
		log.Println("GetLastDocIndexByType", "ehrId", ehrId, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Getting doc from storage
	docDecrypted, err := h.Doc.GetDocFromStorageById(userId, doc.StorageId, []byte(ehrId))
	if err != nil {
		//TODO some logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document getting from storage error"})
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

// GetBySubjectIdAndNamespace
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
func (h EhrHandler) GetBySubjectIdAndNamespace(c *gin.Context) {
	subjectId := c.Query("subject_id")
	namespace := c.Query("subject_namespace")
	if subjectId == "" || namespace == "" {
		log.Println("Subject data is not filled correctly")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docDecrypted, err := h.GetDocBySubject(userId, subjectId, namespace)
	if err != nil {
		log.Println("Can't get document by subject", err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

func (h *EhrHandler) respondWithDocOrHeaders(doc *model.EHR, c *gin.Context) {
	c.Header("Location", h.Cfg.BaseUrl+"/v1/ehr/"+doc.EhrId.Value)
	c.Header("ETag", doc.EhrId.Value)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
