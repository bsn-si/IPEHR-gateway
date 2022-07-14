package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/composition"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type CompositionHandler struct {
	cfg     *config.Config
	service *composition.Service
}

func NewCompositionHandler(docService *service.DefaultDocumentService, cfg *config.Config) *CompositionHandler {
	return &CompositionHandler{
		cfg:     cfg,
		service: composition.NewCompositionService(docService, cfg),
	}
}

// Create
// @Summary      Create COMPOSITION
// @Description  Work in progress...
// @Description  Creates the first version of a new COMPOSITION in the EHR identified by ehr_id.
// @Description
// @Tags     COMPOSITION
// @Accept   json
// @Produce  json
// @Param    ehr_id         path      string                 true   "EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param    AuthUserId     header    string                 true   "UserId - UUID"
// @Param    GroupAccessId  header    string                 false  "GroupAccessId - UUID. If not specified, the default access group will be used."
// @Param    Prefer         header    string                 true   "The new EHR resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param    Request        body      model.SwagComposition  true   "COMPOSITION"
// @Success  201            {object}  model.SwagComposition
// @Header   201            {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/composition/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
// @Header   201            {string}  ETag      "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
// @Failure  400            "Is returned when the request has invalid ehr_id or invalid content (e.g. content could not be converted to a valid COMPOSITION object)"
// @Failure  404            "Is returned when an EHR with ehr_id does not exist."
// @Failure  422            "Is returned when the content could be converted to a COMPOSITION, but there are semantic validation errors, such as the underlying template is not known or is not validating the supplied COMPOSITION)."
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /ehr/{ehr_id}/composition [post]
func (h *CompositionHandler) Create(c *gin.Context) {
	ehrID := c.Param("ehrid")
	if !h.service.Doc.ValidateID(ehrID, types.Ehr) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Request body error"})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic("Cant close body request")
		}
	}(c.Request.Body)

	var request model.Composition

	if err = json.Unmarshal(data, &request); err != nil {
		log.Println("Composition Create request unmarshal error", err)
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

	// TODO move this checks into common???
	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	groupIDStr := c.GetHeader("GroupAccessId")
	if groupIDStr == "" {
		groupIDStr = h.cfg.DefaultGroupAccessID
	}

	groupAccessUUID, err := uuid.Parse(groupIDStr)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	// Composition document creating
	doc, err := h.service.Create(userID, &ehrUUID, &groupAccessUUID, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Composition creating error"})
		return
	}

	h.respondWithDocOrHeaders(ehrID, doc, c)
}

// GetByID
// @Summary      Get COMPOSITION by version id
// @Description  Retrieves a particular version of the COMPOSITION identified by `version_uid` and associated with the EHR identified by `ehr_id`.
// @Description
// @Tags     COMPOSITION
// @Accept   json
// @Produce  json
// @Param    ehr_id       path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param    version_uid  path      string  true  "VERSION identifier taken from VERSION.uid.value. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1"
// @Param    AuthUserId   header    string  true  "UserId UUID"
// @Success  200          {object}  model.SwagComposition
// @Failure  204          "Is returned when the COMPOSITION is deleted (logically)."
// @Failure  400          "Is returned when AuthUserId is not specified"
// @Failure  404          "is returned when an EHR with `ehr_id` does not exist or when an COMPOSITION with `version_uid` does not exist."
// @Failure  500          "Is returned when an unexpected error occurs while processing a request"
// @Router   /ehr/{ehr_id}/composition/{version_uid} [get]
func (h *CompositionHandler) GetByID(c *gin.Context) {
	ehrID := c.Param("ehrid")
	if !h.service.Doc.ValidateID(ehrID, types.Ehr) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("version_uid")
	if !h.service.Doc.ValidateID(versionUID, types.Composition) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.service.GetByID(userID, &ehrUUID, versionUID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrAlreadyDeleted) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, data)
}

// Delete
// @Summary      Deletes the COMPOSITION by version id
// @Description  Deletes the COMPOSITION identified by `preceding_version_uid` and associated with the EHR identified by `ehr_id`.
// @Description
// @Tags         COMPOSITION
// @Accept       json
// @Produce      json
// @Param        ehr_id       path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        preceding_version_uid  path      string  true  "Identifier of the COMPOSITION to be deleted. This MUST be the last (most recent) version. Example: `8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1`"
// @Param        AuthUserId   header    string  true  "UserId UUID"
// @Failure      204          "`No Content` is returned when COMPOSITION was deleted."
// @Failure      400          "`Bad Request` is returned when the composition with `preceding_version_uid` is already deleted."
// @Failure      404          "`Not Found` is returned when an EHR with ehr_id does not exist or when a COMPOSITION with preceding_version_uid does not exist."
// @Failure      409          "`Conflict` is returned when supplied `preceding_version_uid` doesn’t match the latest version. Returns latest version in the Location and ETag headers."
// @Failure      500          "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/composition/{preceding_version_uid} [delete]
func (h *CompositionHandler) Delete(c *gin.Context) {
	ehrID := c.Param("ehrid")
	if !h.service.Doc.ValidateID(ehrID, types.Ehr) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("preceding_version_uid")
	if !h.service.Doc.ValidateID(versionUID, types.Composition) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	newUID, err := h.service.DeleteByID(userID, &ehrUUID, versionUID)

	switch {
	case err == nil:
		h.addResponseHeaders(ehrID, newUID, c)
		c.AbortWithStatus(http.StatusNoContent)
	case errors.Is(err, errors.ErrAlreadyDeleted):
		c.AbortWithStatus(http.StatusBadRequest)
	case errors.Is(err, errors.ErrIsNotExist):
		c.AbortWithStatus(http.StatusNotFound)
	default:
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

// Update
// @Summary      Updates the COMPOSITION by version id
// @Description  Updates COMPOSITION identified by `versioned_object_uid` and associated with the EHR
// @Description  identified by `ehr_id`. If the request body already contains a COMPOSITION.uid.value,
// @Description  it must match the `versioned_object_uid` in the URL. The existing latest `version_uid`
// @Description  of COMPOSITION resource (i.e the `preceding_version_uid`) must be specified in the `If-Match` header.
// @Description
// @Tags         COMPOSITION
// @Accept       json
// @Produce      json
// @Param        ehr_id       path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        versioned_object_uid  path      string  true  "identifier of the COMPOSITION to be updated. Example: `8849182c-82ad-4088-a07f-48ead4180515`"
// @Param        AuthUserId ::openEHRSys.example.com::1  header    string  true  "UserId UUID"
// @Param        Prefer      header    string                 true  "The updated COMPOSITION resource is returned to the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        If-Match    header    string                 true  "The existing latest version_uid of COMPOSITION resource (i.e the preceding_version_uid). Example: `8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1`"
// @Param        Request     body      model.SwagComposition  true  "List of changes in COMPOSITION"
// @Success      200         {object}  model.SwagComposition  true  "Is returned when the COMPOSITION is successfully updated and the updated resource is returned in the body when Prefer header value is `return=representation.`"
// @Header       200         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/composition/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Header       200         {string}  ETag      "8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Failure      422          "`Unprocessable Entity` is returned when the content could be converted to a COMPOSITION, but there are semantic validation errors, such as the underlying template is not known or is not validating the supplied COMPOSITION)."
// @Failure      400          "`Bad Request` is returned when the request has invalid `ehr_id` or invalid content (e.g. either the body of the request could not be read, or converted to a valid COMPOSITION object)"
// @Failure      404          "`Not Found` is returned when an EHR with ehr_id does not exist or when a COMPOSITION with version_object_uid does not exist."
// @Failure      412          "`Version conflict` is returned when `If-Match` request header doesn’t match the latest version (of this versioned object) on the service side. Returns also latest `version_uid` in the `Location` and `ETag` headers."
// @Failure      500          "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/composition/{versioned_object_uid} [put]
func (h CompositionHandler) Update(c *gin.Context) {
	ehrID := c.Param("ehrid")
	if !h.service.Doc.ValidateID(ehrID, types.Ehr) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUID := c.Param("versioned_object_uid")
	if !h.service.Doc.ValidateID(versionUID, types.Composition) {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	precedingVersionUID := c.GetHeader("If-Match")
	if precedingVersionUID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "If-Match is empty"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Checking EHR does not exist
	_, err = h.service.Doc.EhrsIndex.Get(userID)
	if errors.Is(err, errors.ErrIsNotExist) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	groupIDStr := c.GetHeader("GroupAccessId")
	if groupIDStr == "" {
		groupIDStr = h.cfg.DefaultGroupAccessID
	}

	groupAccessUUID, err := uuid.Parse(groupIDStr)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var compositionUpdate model.Composition

	if err = json.Unmarshal(data, &compositionUpdate); err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if !compositionUpdate.Validate() {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	if compositionUpdate.UID.Value != precedingVersionUID {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	compositionLast, err := h.service.GetLastByBaseID(userID, &ehrUUID, versionUID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	if compositionLast.UID.Value != precedingVersionUID {
		h.addResponseHeaders(ehrID, compositionLast.UID.Value, c)
		c.AbortWithStatus(http.StatusPreconditionFailed)

		return
	}

	compositionUpdated, err := h.service.Update(userID, &ehrUUID, &groupAccessUUID, &compositionUpdate)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.addResponseHeaders(ehrID, compositionUpdated.UID.Value, c)
	c.JSON(http.StatusOK, compositionUpdated)
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
	c.Header("Location", h.cfg.BaseURL+"/v1/ehr/"+ehrID+"/composition/"+uid)
	c.Header("ETag", uid)
}
