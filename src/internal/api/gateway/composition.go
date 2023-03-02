package gateway

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/composition"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
)

type (
	CompositionService interface {
		helper.Finder
		DefaultGroupAccess() *uuid.UUID
		Create(ctx context.Context, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition, procRequest *proc.Request) (*model.Composition, error)
		Update(ctx context.Context, procRequest *proc.Request, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition) (*model.Composition, error)
		GetLastByBaseID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, versionUID string) (*model.Composition, error)
		GetByID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, versionUID string) (*model.Composition, error)
		DeleteByID(ctx context.Context, procRequest *proc.Request, ehrUUID *uuid.UUID, versionUID, userID, systemID string) (string, error)
		GetList(ctx context.Context, userID, systemID string) ([]*model.EhrDocumentItem, error)
	}

	Indexer interface {
		GetEhrUUIDByUserID(ctx context.Context, userID, systemID string) (*uuid.UUID, error)
	}

	ProcessingService interface {
		NewRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (*processing.Request, error)
	}

	CompositionHandler struct {
		service       CompositionService
		indexer       Indexer
		processingSvc ProcessingService
		baseURL       string
	}
)

func NewCompositionHandler(docService *service.DefaultDocumentService, compositionService *composition.Service, baseURL string) *CompositionHandler {
	return &CompositionHandler{
		service:       compositionService,
		indexer:       docService.Infra.Index,
		processingSvc: docService.Proc,
		baseURL:       baseURL,
	}
}

// Create
//
//	@Summary		Create COMPOSITION
//	@Description	Work in progress...
//	@Description	Creates the first version of a new COMPOSITION in the EHR identified by ehr_id.
//	@Description
//	@Tags		COMPOSITION
//	@Accept		json
//	@Produce	json
//	@Param		ehr_id			path		string				true	"EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		Authorization	header		string				true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string				true	"UserId"
//	@Param		EhrSystemId		header		string				false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		GroupAccessId	header		string				false	"GroupAccessId - UUID. If not specified, the default access group will be used."
//	@Param		Prefer			header		string				true	"The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
//	@Param		Request			body		model.Composition	true	"COMPOSITION"
//	@Success	201				{object}	model.Composition
//	@Header		201				{string}	Location	"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/composition/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
//	@Header		201				{string}	ETag		"8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
//	@Header		201				{string}	RequestID	"Request identifier"
//	@Failure	400				"Is returned when the request has invalid ehr_id or invalid content (e.g. content could not be converted to a valid COMPOSITION object)"
//	@Failure	404				"Is returned when an EHR with ehr_id does not exist."
//	@Failure	422				"Is returned when the content could be converted to a COMPOSITION, but there are semantic validation errors, such as the underlying template is not known or is not validating the supplied COMPOSITION)."
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/ehr/{ehr_id}/composition [post]
func (h *CompositionHandler) Create(c *gin.Context) {
	systemID := c.GetString("ehrSystemID")
	reqID := c.GetString("reqID")
	ehrID := c.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	userEhrUUID, err := h.indexer.GetEhrUUIDByUserID(c, userID, systemID)
	switch {
	case err != nil && errors.Is(err, errors.ErrIsNotExist):
		c.AbortWithStatus(http.StatusNotFound)
		return
	case err != nil:
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userEhrUUID.String() != ehrUUID.String() {
		log.Printf("userEhrUUID and ehrUUID is not equal: %s != %s", userEhrUUID, ehrUUID)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	groupAccessUUID := h.service.DefaultGroupAccess()

	if c.GetHeader("GroupAccessId") != "" {
		UUID, err := uuid.Parse(c.GetHeader("GroupAccessId"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "GroupAccessId parsing error"})
			return
		}

		groupAccessUUID = &UUID
	}

	composition := &model.Composition{}

	if err := json.NewDecoder(c.Request.Body).Decode(composition); err != nil {
		log.Println("Composition Create request unmarshal error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body parsing error"})
		return
	}
	defer c.Request.Body.Close()

	if ok, _ := composition.Validate(); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	procRequest, err := h.processingSvc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestCompositionCreate)
	if err != nil {
		log.Println("Composition create NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Composition document creating
	doc, err := h.service.Create(c, userID, systemID, &ehrUUID, groupAccessUUID, composition, procRequest)
	if err != nil {
		log.Println("Composition creating error:", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Composition creating error"})
		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("Composition procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.respondWithDocOrHeaders(ehrID, doc, c)
}

// GetByID
//
//	@Summary		Get COMPOSITION by version id
//	@Description	Retrieves a particular version of the COMPOSITION identified by `version_uid` and associated with the EHR identified by `ehr_id`.
//	@Description
//	@Tags		COMPOSITION
//	@Accept		json
//	@Produce	json
//	@Param		ehr_id			path		string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		version_uid		path		string	true	"VERSION identifier taken from VERSION.uid.value. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
//	@Param		Authorization	header		string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string	true	"UserId"
//	@Param		EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success	200				{object}	model.Composition
//	@Success	202				"Is returned when the request is still being processed"
//	@Failure	204				"Is returned when the COMPOSITION is deleted (logically)."
//	@Failure	400				"Is returned when AuthUserId is not specified"
//	@Failure	404				"is returned when an EHR with `ehr_id` does not exist or when an COMPOSITION with `version_uid` does not exist."
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/ehr/{ehr_id}/composition/{version_uid} [get]
func (h *CompositionHandler) GetByID(c *gin.Context) {
	ehrID := c.Param("ehrid")

	systemID := c.GetString("ehrSystemID")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("version_uid")

	//TODO validate versionUID

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	// Checking EHR does not exist
	userEhrUUID, err := h.indexer.GetEhrUUIDByUserID(c, userID, systemID)
	switch {
	case err != nil && errors.Is(err, errors.ErrIsNotExist):
		c.AbortWithStatus(http.StatusNotFound)
		return
	case err != nil:
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userEhrUUID.String() != ehrUUID.String() {
		log.Printf("userEhrUUID and ehrUUID is not equal: %s != %s", userEhrUUID, ehrUUID)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	data, err := h.service.GetByID(c, userID, systemID, &ehrUUID, versionUID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else if errors.Is(err, errors.ErrAlreadyDeleted) {
			c.AbortWithStatus(http.StatusNoContent)
		} else if errors.Is(err, errors.ErrIsInProcessing) {
			c.AbortWithStatus(http.StatusAccepted)
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, data)
}

// Delete
//
//	@Summary		Deletes the COMPOSITION by version id
//	@Description	Deletes the COMPOSITION identified by `preceding_version_uid` and associated with the EHR identified by `ehr_id`.
//	@Description
//	@Tags		COMPOSITION
//	@Accept		json
//	@Produce	json
//	@Param		ehr_id					path	string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		preceding_version_uid	path	string	true	"Identifier of the COMPOSITION to be deleted. This MUST be the last (most recent)  version.  Example:  `8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1`"
//	@Param		Authorization			header	string	true	"Bearer AccessToken"
//	@Param		AuthUserId				header	string	true	"UserId"
//	@Param		EhrSystemId				header	string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Failure	204						"`No Content` is returned when COMPOSITION was deleted."
//	@Header		204						{string}	RequestID	"Request identifier"
//	@Failure	400						"`Bad Request` is returned when the composition with `preceding_version_uid` is already deleted."
//	@Failure	404						"`Not Found` is returned when an EHR with ehr_id does not exist or when a COMPOSITION with preceding_version_uid does not exist."
//	@Failure	409						"`Conflict` is returned when supplied `preceding_version_uid` doesn’t match the latest version. Returns latest version in the Location and ETag headers."
//	@Failure	500						"Is returned when an unexpected error occurs while processing a request"
//	@Router		/ehr/{ehr_id}/composition/{preceding_version_uid} [delete]
func (h *CompositionHandler) Delete(c *gin.Context) {
	ehrID := c.Param("ehrid")

	systemID := c.GetString("ehrSystemID")

	//TODO validate ehrID

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("preceding_version_uid")

	//TODO validate versionUID

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	userEhrUUID, err := h.indexer.GetEhrUUIDByUserID(c, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	if userEhrUUID.String() != ehrUUID.String() {
		log.Printf("userEhrUUID and ehrUUID is not equal: %s != %s", userEhrUUID, ehrUUID)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	reqID := c.GetString("reqID")

	procRequest, err := h.processingSvc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestCompositionDelete)
	if err != nil {
		log.Println("Composition delete NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	newUID, err := h.service.DeleteByID(c, procRequest, &ehrUUID, versionUID, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		} else if errors.Is(err, errors.ErrAlreadyDeleted) {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			log.Println("DeleteByID error:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("Composition delete procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.addResponseHeaders(ehrID, newUID, c)
	c.AbortWithStatus(http.StatusNoContent)
}

// Update
//
//	@Summary		Updates the COMPOSITION by version id
//	@Description	Updates COMPOSITION identified by `versioned_object_uid` and associated with the EHR
//	@Description	identified by `ehr_id`. If the request body already contains a COMPOSITION.uid.value,
//	@Description	it must match the `versioned_object_uid` in the URL. The existing latest `version_uid`
//	@Description	of COMPOSITION resource (i.e the `preceding_version_uid`) must be specified in the `If-Match` header.
//	@Description
//	@Tags		COMPOSITION
//	@Accept		json
//	@Produce	json
//	@Param		ehr_id					path		string				true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		versioned_object_uid	path		string				true	"identifier of the COMPOSITION to be updated. Example: `8849182c-82ad-4088-a07f-48ead4180515`"
//	@Param		Authorization			header		string				true	"Bearer AccessToken"
//	@Param		AuthUserId				header		string				true	"UserId"
//	@Param		EhrSystemId				header		string				false	"The identifier of the system, typically a reverse domain identifier"
//	@Param		Prefer					header		string				true	"The updated COMPOSITION resource is returned to the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
//	@Param		If-Match				header		string				true	"The existing latest version_uid of COMPOSITION resource (i.e the preceding_version_uid).  Example:  `8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1`"
//	@Param		Request					body		model.Composition	true	"List of changes in COMPOSITION"
//	@Success	200						{object}	model.Composition	"Is returned when the COMPOSITION is successfully updated and the updated resource is returned in the body when Prefer header value is `return=representation.`"
//	@Header		200						{string}	Location			"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/composition/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Header		200						{string}	ETag				"8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Header		200						{string}	RequestID			"Request identifier"
//	@Failure	422						"`Unprocessable Entity` is returned when the content could be converted to a COMPOSITION, but there are semantic validation errors, such as the underlying template is not known or is not validating the supplied COMPOSITION)."
//	@Failure	400						"`Bad Request` is returned when the request has invalid `ehr_id` or invalid content (e.g. either the body of the request could not be read, or converted to a valid COMPOSITION object)"
//	@Failure	404						"`Not Found` is returned when an EHR with ehr_id does not exist or when a COMPOSITION with version_object_uid does not exist."
//	@Failure	412						"`Version conflict` is returned when `If-Match` request header doesn’t match the latest version (of this versioned object)  on  the  service  side.  Returns  also  latest  `version_uid`  in  the  `Location`  and  `ETag`  headers."
//	@Failure	500						"Is returned when an unexpected error occurs while processing a request"
//	@Router		/ehr/{ehr_id}/composition/{versioned_object_uid} [put]
func (h CompositionHandler) Update(c *gin.Context) {
	systemID := c.GetString("ehrSystemID")
	reqID := c.GetString("reqID")
	ehrID := c.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("versioned_object_uid")

	//TODO validate versionUID

	precedingVersionUID := c.GetHeader("If-Match")
	if precedingVersionUID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "If-Match is empty"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	userEhrUUID, err := h.indexer.GetEhrUUIDByUserID(c, userID, systemID)
	switch {
	case err != nil && errors.Is(err, errors.ErrIsNotExist):
		c.AbortWithStatus(http.StatusNotFound)
		return
	case err != nil:
		log.Println("GetEhrIDByUser error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userEhrUUID.String() != ehrUUID.String() {
		log.Printf("userEhrUUID and ehrUUID is not equal: %s != %s", userEhrUUID, ehrUUID)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	groupAccessUUID := h.service.DefaultGroupAccess()

	if c.GetHeader("GroupAccessId") != "" {
		UUID, err := uuid.Parse(c.GetHeader("GroupAccessId"))
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "GroupAccessId parsing error"})
			return
		}

		groupAccessUUID = &UUID
	}

	compositionUpdate := model.Composition{}

	if err := json.NewDecoder(c.Request.Body).Decode(&compositionUpdate); err != nil {
		log.Println("Composition Update request unmarshal error", err)
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Request body parsing error"})
		return
	}
	defer c.Request.Body.Close()

	if ok, _ := compositionUpdate.Validate(); !ok {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if compositionUpdate.UID.Value != precedingVersionUID {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	compositionLast, err := h.service.GetLastByBaseID(c, userID, systemID, &ehrUUID, versionUID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrIsInProcessing) {
			c.AbortWithStatus(http.StatusAccepted)
			return
		}

		log.Println("GetLastByBaseID error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	if compositionLast.UID.Value != precedingVersionUID {
		h.addResponseHeaders(ehrID, compositionLast.UID.Value, c)
		c.AbortWithStatus(http.StatusPreconditionFailed)

		return
	}

	procRequest, err := h.processingSvc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestCompositionUpdate)
	if err != nil {
		log.Println("Compocition update NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	compositionUpdated, err := h.service.Update(c, procRequest, userID, systemID, &ehrUUID, groupAccessUUID, &compositionUpdate)
	if err != nil {
		log.Println("Composition Update error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("Composition update procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.addResponseHeaders(ehrID, compositionUpdated.UID.Value, c)
	c.JSON(http.StatusOK, compositionUpdated)
}

// List
//
//	@Summary		Get all COMPOSITIONs
//	@Description	Retrieves all versions of all COMPOSITIONs associated with the EHR identified by `ehr_id`.
//	@Description
//	@Tags		COMPOSITION
//	@Accept		json
//	@Produce	json
//	@Param		ehr_id			path		string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param		Authorization	header		string	true	"Bearer AccessToken"
//	@Param		AuthUserId		header		string	true	"UserId"
//	@Param		EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success	200				{object}	[]model.EhrDocumentItem
//	@Failure	400				"Is returned when AuthUserId or EhrSystemId is not specified"
//	@Failure	404				"is returned when an EHR with `ehr_id` does not exist."
//	@Failure	500				"Is returned when an unexpected error occurs while processing a request"
//	@Router		/ehr/{ehr_id}/composition [get]
func (h CompositionHandler) GetList(c *gin.Context) {
	systemID := c.GetString("ehrSystemID")
	userID := c.GetString("userID")

	list, err := h.service.GetList(c, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
		}

		log.Println("Composition GetList error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	if list == nil {
		list = []*model.EhrDocumentItem{}
	}

	c.JSON(http.StatusOK, list)
}

func (h *CompositionHandler) respondWithDocOrHeaders(ehrID string, doc *model.Composition, c *gin.Context) {
	uid := doc.UID.Value
	h.addResponseHeaders(ehrID, uid, c)

	prefer := c.Request.Header.Get("Prefer")
	//nolint:goconst
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}

func (h *CompositionHandler) addResponseHeaders(ehrID string, uid string, c *gin.Context) {
	c.Header("Location", h.baseURL+"/v1/ehr/"+ehrID+"/composition/"+uid)
	c.Header("ETag", uid)
}
