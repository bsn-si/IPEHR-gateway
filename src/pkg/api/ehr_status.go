package api

import (
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr_status"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EhrStatusHandler struct {
	service *ehr_status.EhrStatusService
}

func NewEhrStatusHandler(docService *service.DefaultDocumentService) *EhrStatusHandler {
	return &EhrStatusHandler{
		service: ehr_status.NewEhrStatusService(docService),
	}
}

func (h EhrStatusHandler) Update(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.service.DocService.ValidateId(ehrId, types.EHR) == false {
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
	docIndexLast, err := h.service.DocService.GetLastDocIndexByType(ehrId, types.EHR_STATUS)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			c.Writer.Header().Set("Location", "") //TODO
			c.Writer.Header().Set("ETag", "")
			c.AbortWithStatus(http.StatusPreconditionFailed)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Find previous EHR_STATUS error: " + err.Error()})
		}
		return
	}

	// Getting doc from storage
	docBytes, err := h.service.DocService.GetDocFromStorageById(userId, docIndexLast.StorageId, []byte(IfMatch))
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
	//TODO baseUrl
	c.Writer.Header().Set("Location", "{baseUrl}/ehr/"+ehrId+"/ehr_status/"+docLast.Uid.Value)
	c.Writer.Header().Set("ETag", docLast.Uid.Value)

	// Checking If-Match header
	if IfMatch != docLast.Uid.Value {
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

	c.Data(http.StatusOK, "application/json", data)
}

func (h EhrStatusHandler) GetById(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.service.DocService.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	versionUid := c.Param("versionid")
	if h.service.DocService.ValidateId(versionUid, types.EHR_STATUS) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	docIndex, err := h.service.DocService.GetDocIndexByDocId(userId, ehrId, versionUid, types.EHR_STATUS)
	if err != nil {
		log.Printf("GetDocIndexByDocId userId: %s ehrId: %s versionId: %s error: %v", userId, ehrId, versionUid, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := h.service.DocService.GetDocFromStorageById(userId, docIndex.StorageId, []byte(versionUid))
	if err != nil {
		//TODO logging
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
