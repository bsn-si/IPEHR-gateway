package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
)

type EhrStatusHandler struct {
	ehrStatusService *ehr.EhrStatusService
}

func NewEhrStatusHandler(docService *service.DefaultDocumentService) *EhrStatusHandler {
	return &EhrStatusHandler{
		ehrStatusService: ehr.NewEhrStatusService(docService),
	}
}

func (h EhrStatusHandler) GetStatus(c *gin.Context) {
	if ok := c.Query("version_at_time"); ok != "" {
		h.GetStatusByTime(c)
		return
	}
	c.AbortWithStatus(http.StatusBadRequest)
}

func (h EhrStatusHandler) Update(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.ehrStatusService.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	IfMatch := c.Request.Header.Get("If-Match")
	docLast, err := h.ehrStatusService.Get(userId, ehrId)
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
		setLocationAndETagHeaders(ehrId, docLast.Uid.Value, c)
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

	update, err := h.ehrStatusService.ParseJson(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	if !h.ehrStatusService.Validate(update) {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err = h.ehrStatusService.Save(ehrId, userId, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR_STATUS saving error"})
		return
	}

	setLocationAndETagHeaders(ehrId, update.Uid.Value, c)

	switch c.Request.Header.Get("Prefer") {
	case "return=representation":
		c.Data(http.StatusOK, "application/json", data)
	case "return=minimal":
		fallthrough
	default:
		c.JSON(http.StatusNoContent, nil)
	}
}

func (h EhrStatusHandler) GetStatusByTime(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.ehrStatusService.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.ehrStatusService.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
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

	docIndex, err := h.ehrStatusService.Doc.GetDocIndexByNearestTime(ehrId, statusTime, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByNearestTime: ehrId: %s statusTime: %s error: %v", ehrId, statusTime, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.ehrStatusService.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h EhrStatusHandler) GetById(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.ehrStatusService.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.ehrStatusService.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docIndex, err := h.ehrStatusService.Doc.GetDocIndexByDocId(userId, ehrId, versionUid, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByDocId userId: %s ehrId: %s versionId: %s error: %v", userId, ehrId, versionUid, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.ehrStatusService.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func setLocationAndETagHeaders(ehrId string, ehrStatusId string, c *gin.Context) {
	c.Header("Location", AppConfig.BaseUrl+"/ehr/"+ehrId+"/ehr_status/"+ehrStatusId)
	c.Header("ETag", ehrStatusId)
}
