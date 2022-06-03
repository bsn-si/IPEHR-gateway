package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type EhrHandler struct {
	service *ehr.EhrService
}

func NewEhrHandler(docService *service.DefaultDocumentService) *EhrHandler {
	return &EhrHandler{
		service: ehr.NewEhrService(docService),
	}
}

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
	_, err = h.service.Doc.EhrsIndex.Get(userId)
	if !errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// EHR document creating
	doc, err := h.service.Create(userId, &request)
	if err != nil {
		if errors.Is(err, errors.AlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}
		return
	}

	respondWithDocOrHeaders(doc, c)
}

func (h EhrHandler) CreateWithId(c *gin.Context) {
	ehrId := c.Param("ehrid")

	// Checking EHR does not exist
	doc, err := h.service.Doc.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil && !errors.Is(err, errors.IsNotExist) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "EHR retrieve error"})
		return
	}
	if doc != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		return
	}

	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
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
	_, err = h.service.Doc.EhrsIndex.Get(userId)
	if !errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// EHR document creating
	newDoc, err := h.service.CreateWithId(userId, ehrId, &request)
	if err != nil {
		if errors.Is(err, errors.AlreadyExist) {
			c.JSON(http.StatusConflict, gin.H{"error": "EHR already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "EHR creating error"})
		}
		return
	}

	respondWithDocOrHeaders(newDoc, c)
}

func (h EhrHandler) GetById(c *gin.Context) {
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

	// Getting docStorageId
	doc, err := h.service.Doc.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil {
		log.Println("GetLastDocIndexByType", "ehrId", ehrId, err)
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Getting doc from storage
	docDecrypted, err := h.service.Doc.GetDocFromStorageById(userId, doc.StorageId, []byte(ehrId))
	if err != nil {
		//TODO some logging
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Document getting from storage error"})
		return
	}

	c.Data(http.StatusOK, "application/json", docDecrypted)
}

func respondWithDocOrHeaders(doc *model.EHR, c *gin.Context) {
	c.Header("Location", AppConfig.BaseUrl+"/v1/ehr/"+doc.EhrId.Value)
	c.Header("ETag", doc.EhrId.Value)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
