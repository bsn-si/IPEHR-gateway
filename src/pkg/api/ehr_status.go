package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
)

type EhrStatusHandler struct {
	*ehr.EhrStatusService
}

func NewEhrStatusHandler(docService *service.DefaultDocumentService, cfg *config.Config) *EhrStatusHandler {
	return &EhrStatusHandler{
		ehr.NewEhrStatusService(docService, cfg),
	}
}

func (h EhrStatusHandler) Get(c *gin.Context) {
	if ok := c.Query("version_at_time"); ok != "" {
		h.GetStatusByTime(c)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

// Update
// @Summary      Update EHR_STATUS
// @Description  Updates EHR_STATUS associated with the EHR identified by `ehr_id`. The existing latest `version_uid` of EHR_STATUS resource (i.e. the `preceding_version_uid`) must be specified in the `If-Match` header. The response will contain the updated EHR_STATUS resource when the `Prefer` header has a value of `return=representation`
// @Tags         EHR_STATUS
// @Accept       json
// @Produce      json
// @Param        ehr_id      path      string                 true  "EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        AuthUserId  header    string                 true  "UserId UUID"
// @Param        Prefer      header    string                 true  "Updated resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
// @Param        Request     body      model.EhrStatusUpdate  true  "EHR_STATUS"
// @Success      200         {object}  model.EhrStatusUpdate
// @Header       200         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/ehr_status/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Header       200         {string}  ETag      "uid of created document. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Success      204         "Is returned when `Prefer` header is missing or is set to `return=minimal`"
// @Header       204         {string}  Location  "{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/ehr_status/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Header       204         {string}  ETag      "uid of created document. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Failure      400         "Is returned when the request has invalid content."
// @Failure      404         "Is returned when an EHR with ehr_id does not exist."
// @Failure      412         "Is returned when `If-Match` request header doesn’t match the latest version on the service side. Returns also latest `version_uid` in the `Location` and `ETag` headers."
// @Failure      500          "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/ehr_status [put]
func (h EhrStatusHandler) Update(c *gin.Context) {
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

	IfMatch := c.Request.Header.Get("If-Match")
	docLast, err := h.GetStatus(userId, ehrId)
	if err != nil {
		log.Println("ehrStatusService.Get error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting last EHR document status"})
	}

	if docLast.Uid == nil || docLast.Uid.Value == "" {
		c.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	// Checking If-Match header
	if IfMatch != docLast.Uid.Value {
		h.setLocationAndETagHeaders(ehrId, docLast.Uid.Value, c)
		c.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	err = c.Request.Body.Close()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	update, err := h.ParseJson(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	if !h.Validate(update) {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err = h.SaveStatus(ehrId, userId, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR_STATUS saving error"})
		return
	}

	h.setLocationAndETagHeaders(ehrId, update.Uid.Value, c)

	switch c.Request.Header.Get("Prefer") {
	case "return=representation":
		c.Data(http.StatusOK, "application/json", data)
	case "return=minimal":
		fallthrough
	default:
		c.JSON(http.StatusNoContent, nil)
	}
}

// GetStatusByTime
// @Summary      Get EHR_STATUS version by time
// @Description  Retrieves a version of the EHR_STATUS associated with the EHR identified by `ehr_id`. If `version_at_time` is supplied, retrieves the version extant at specified time, otherwise retrieves the latest EHR_STATUS version.
// @Tags         EHR_STATUS
// @Accept       json
// @Produce      json
// @Param        ehr_id       path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        version_at_time  query     string  true  "A given time in the extended ISO 8601 format. Example: 2015-01-20T19:30:22.765+01:00"
// @Param        AuthUserId   header    string  true  "UserId UUID"
// @Success      200          {object}  model.EhrStatusUpdate
// @Failure      400              "Is returned when the request has invalid content such as an invalid `version_at_time` format."
// @Failure      404              "Is returned when EHR with `ehr_id` does not exist or a version of an EHR_STATUS resource does not exist at the specified `version_at_time`"
// @Failure      500              "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/ehr_status [get]
func (h EhrStatusHandler) GetStatusByTime(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}
	versionAtTime := c.Query("version_at_time")
	statusTime, err := time.Parse(common.OPENEHR_TIME_FORMAT, versionAtTime)
	if err != nil {
		log.Printf("Incorrect format of time option, use %s", common.OPENEHR_TIME_FORMAT)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	docIndex, err := h.Doc.GetDocIndexByNearestTime(ehrId, statusTime, types.EHR_STATUS)
	if err != nil {
		// TODO: It is necessary to write a log if document not found here?
		//log.Printf("GetDocIndexByNearestTime: ehrId: %s statusTime: %s error: %v", ehrId, statusTime, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

// GetById
// @Summary      Get EHR_STATUS by version id
// @Description  Retrieves a particular version of the EHR_STATUS identified by `version_uid` and associated with the EHR identified by `ehr_id`.
// @Tags         EHR_STATUS
// @Accept       json
// @Produce      json
// @Param        ehr_id           path      string  true  "EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        version_uid  path      string  true  "VERSION identifier taken from VERSION.uid.value. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
// @Param        AuthUserId       header    string  true  "UserId UUID"
// @Success      200              {object}  model.EhrStatusUpdate
// @Failure      400          "Is returned when AuthUserId is not specified"
// @Failure      404          "is returned when an EHR with `ehr_id` does not exist or when an EHR_STATUS with `version_uid` does not exist."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/ehr_status/{version_uid} [get]
func (h EhrStatusHandler) GetById(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docIndex, err := h.Doc.GetDocIndexByDocId(userId, ehrId, versionUid, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByDocId userId: %s ehrId: %s versionId: %s error: %v", userId, ehrId, versionUid, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h *EhrStatusHandler) setLocationAndETagHeaders(ehrId string, ehrStatusId string, c *gin.Context) {
	c.Header("Location", h.Cfg.BaseUrl+"/ehr/"+ehrId+"/ehr_status/"+ehrStatusId)
	c.Header("ETag", ehrStatusId)
}
