package gateway

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	docGroupService "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/docGroup"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/ehr"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/groupAccess"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	userService "github.com/bsn-si/IPEHR-gateway/src/pkg/user/service"
)

type EhrStatusHandler struct {
	service *ehr.Service
	baseURL string
}

func NewEhrStatusHandler(docSvc *service.DefaultDocumentService, userSvc *userService.Service, docGroupSvc *docGroupService.Service, gaSvc *groupAccess.Service, baseURL string) *EhrStatusHandler {
	return &EhrStatusHandler{
		service: ehr.NewService(docSvc, userSvc, docGroupSvc, gaSvc),
		baseURL: baseURL,
	}
}

// Update
//
//	@Summary		Update EHR_STATUS
//	@Description	Updates EHR_STATUS associated with the EHR identified by `ehr_id`.
//	@Tags			EHR_STATUS
//	@Accept			json
//	@Produce		json
//	@Param			ehr_id			path		string					true	"EHR identifier. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param			Authorization	header		string					true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string					true	"UserId"
//	@Param			EhrSystemId		header		string					false	"The identifier of the system, typically a reverse domain identifier"
//	@Param			If-Match		header		string					true	"The existing latest `version_uid` of EHR_STATUS resource (i.e. the `preceding_version_uid`)  must  be  specified."
//	@Param			Prefer			header		string					true	"Updated resource is returned in the body when the request’s `Prefer` header value is `return=representation`, otherwise only headers are returned."
//	@Param			Request			body		model.EhrStatusUpdate	true	"EHR_STATUS"
//	@Success		200				{object}	model.EhrStatusUpdate
//	@Header			200				{string}	Location	"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/ehr_status/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Header			200				{string}	ETag		"uid of created document. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Header			200				{string}	RequestID	"Request identifier"
//	@Success		204				"Is returned when `Prefer` header is missing or is set to `return=minimal`"
//	@Header			204				{string}	Location	"{baseUrl}/ehr/7d44b88c-4199-4bad-97dc-d78268e01398/ehr_status/8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Header			204				{string}	ETag		"uid of created document. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Failure		400				"Is returned when the request has invalid content."
//	@Failure		404				"Is returned when an EHR with ehr_id does not exist."
//	@Failure		412				"Is returned when `If-Match` request header doesn’t match the latest version on the service side. Returns also latest `version_uid` in the `Location` and `ETag` headers."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr/{ehr_id}/ehr_status [put]
func (h *EhrStatusHandler) Update(c *gin.Context) {
	ehrID := c.Param("ehrid")
	reqID := c.GetString("reqID")

	//TODO validate ehrID

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

	systemID := c.GetString("ehrSystemID")

	IfMatch := c.Request.Header.Get("If-Match")

	docLast, err := h.service.GetStatus(c, userID, systemID, &ehrUUID)
	if err != nil {
		log.Println("ehrStatusService.Get error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting last EHR document status"})

		return
	}

	if docLast.UID == nil || docLast.UID.Value == "" {
		c.AbortWithStatus(http.StatusPreconditionFailed)
		return
	}

	// Checking If-Match header
	if IfMatch != docLast.UID.Value {
		h.setLocationAndETagHeaders(ehrID, docLast.UID.Value, c)
		c.AbortWithStatus(http.StatusPreconditionFailed)

		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}

	if err = c.Request.Body.Close(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	var status model.EhrStatus
	if err = json.Unmarshal(data, &status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	if !h.service.ValidateStatus(&status) {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	procRequest, err := h.service.Doc.Proc.NewRequest(reqID, userID, ehrUUID.String(), proc.RequestEhrStatusUpdate)
	if err != nil {
		log.Println("EhrStatus update NewRequest error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.service.UpdateStatus(c, procRequest, userID, systemID, &ehrUUID, &status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("EhrStatus update procRequest commit error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	h.setLocationAndETagHeaders(ehrID, status.UID.Value, c)

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
//
//	@Summary		Get EHR_STATUS version by time
//	@Description	Retrieves a version of the EHR_STATUS associated with the EHR identified by `ehr_id`. If `version_at_time` is supplied, retrieves the version extant at specified time, otherwise retrieves the latest EHR_STATUS version.
//	@Tags			EHR_STATUS
//	@Accept			json
//	@Produce		json
//	@Param			ehr_id			path		string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param			version_at_time	query		string	true	"A given time in the extended ISO 8601 format. Example: 2015-01-20T19:30:22.765+01:00"
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Param			EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200				{object}	model.EhrStatusUpdate
//	@Success		202				"Is returned when the request is still being processed"
//	@Failure		400				"Is returned when the request has invalid content such as an invalid `version_at_time` format."
//	@Failure		404				"Is returned when EHR with `ehr_id` does not exist or a version of an EHR_STATUS resource does not exist at the specified `version_at_time`"
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr/{ehr_id}/ehr_status [get]
func (h *EhrStatusHandler) GetStatusByTime(c *gin.Context) {
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

	systemID := c.GetString("ehrSystemID")

	versionAtTime := c.Query("version_at_time")

	statusTime, err := time.Parse(common.OpenEhrTimeFormat, versionAtTime)
	if err != nil {
		log.Printf("Incorrect format of time option, use %s", common.OpenEhrTimeFormat)
		c.AbortWithStatus(http.StatusNotFound)

		return
	}

	docDecrypted, err := h.service.GetStatusByNearestTime(c, userID, systemID, &ehrUUID, statusTime)
	if err != nil {
		if errors.Is(err, errors.ErrIsInProcessing) {
			c.AbortWithStatus(http.StatusAccepted)
			return
		} else if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("service.GetStatusByNearestTime error: %v userID: %s ehrID: %s", err, userID, ehrID)
		c.AbortWithStatus(http.StatusInternalServerError)

		return
	}

	c.JSON(http.StatusOK, docDecrypted)
}

// GetByID
//
//	@Summary		Get EHR_STATUS by version id
//	@Description	Retrieves a particular version of the EHR_STATUS identified by `version_uid` and associated with the EHR identified by `ehr_id`.
//	@Tags			EHR_STATUS
//	@Accept			json
//	@Produce		json
//	@Param			ehr_id			path		string	true	"EHR identifier taken from EHR.ehr_id.value. Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
//	@Param			version_uid		path		string	true	"VERSION identifier taken from VERSION.uid.value. Example: 8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::2"
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Param			EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200				{object}	model.EhrStatusUpdate
//	@Success		202				"Is returned when the request is still being processed"
//	@Failure		400				"Is returned when AuthUserId is not specified"
//	@Failure		404				"is returned when an EHR with `ehr_id` does not exist or when an EHR_STATUS with `version_uid` does not exist."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/ehr/{ehr_id}/ehr_status/{version_uid} [get]
func (h *EhrStatusHandler) GetByID(c *gin.Context) {
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

	systemID := c.GetString("ehrSystemID")

	versionUID := c.Param("versionid")

	objectVersionID, err := base.NewObjectVersionID(versionUID, systemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ehrSystemID not match with versionUID"})
		return
	}

	docDecrypted, err := h.service.GetStatusByVersionID(c, userID, systemID, &ehrUUID, objectVersionID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		} else if errors.Is(err, errors.ErrIsInProcessing) {
			c.AbortWithStatus(http.StatusAccepted)
			return
		}

		log.Printf("service.GetStatusByVersionID error: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

func (h *EhrStatusHandler) setLocationAndETagHeaders(ehrID string, ehrStatusID string, c *gin.Context) {
	c.Header("Location", h.baseURL+"/ehr/"+ehrID+"/ehr_status/"+ehrStatusID)
	c.Header("ETag", ehrStatusID)
}
