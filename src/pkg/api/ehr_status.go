package api

import (
	"github.com/gin-gonic/gin"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const ALLOWED_TIME_FORMAT = "2006-01-02T15:04:05.999-07:00"

type EhrStatusHandler struct {
	service *ehr.EhrStatusService
}

func NewEhrStatusHandler(docService *service.DefaultDocumentService) *EhrStatusHandler {
	return &EhrStatusHandler{
		service: ehr.NewEhrStatusService(docService),
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
	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	IfMatch := c.Request.Header.Get("If-Match")

	// Getting previous EHR_STATUS
	docIndexLast, err := h.service.Doc.DocsIndex.GetLastByType(ehrId, types.EHR_STATUS)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			setLocationAndETagHeaders(ehrId, "", c)
			c.AbortWithStatus(http.StatusPreconditionFailed)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Find previous EHR_STATUS error: " + err.Error()})
		}
		return
	}

	// Getting doc from storage
	docBytes, err := h.service.Doc.GetDocFromStorageById(userId, docIndexLast.StorageId, []byte(IfMatch))
	if err != nil {
		log.Println("GetDocFromStorageById error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document getting from storage error"})
		return
	}

	docLast, err := h.service.ParseJson(docBytes)
	if err != nil {
		log.Println("EHR_STATUS from storage parsing error error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document from storage parsing error"})
		return
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
	c.Request.Body.Close()

	update, err := h.service.ParseJson(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request content is invalid"})
		return
	}

	if !h.service.Validate(update) {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	if err = h.service.Save(ehrId, userId, update); err != nil {
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
	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.service.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}
	versionAtTime := c.Query("version_at_time")
	statusTime, err := time.Parse(ALLOWED_TIME_FORMAT, versionAtTime)

	if err != nil {
		log.Printf("Incorrect format of time option, use %s", ALLOWED_TIME_FORMAT)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	docIndex, err := h.service.Doc.GetDocIndexByNearestTime(ehrId, statusTime, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByNearestTime: ehrId: %s statusTime: %s error: %v", ehrId, statusTime, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.service.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

func (h EhrStatusHandler) GetById(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.service.Doc.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docIndex, err := h.service.Doc.GetDocIndexByDocId(userId, ehrId, versionUid, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByDocId userId: %s ehrId: %s versionId: %s error: %v", userId, ehrId, versionUid, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.service.Doc.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
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
