package api

import (
	"encoding/json"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service/composition"
	"hms/gateway/pkg/docs/types"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type CompositionHandler struct {
	service *composition.CompositionService
}

func NewCompositionHandler(docService *service.DefaultDocumentService, cfg *config.Config) *CompositionHandler {
	return &CompositionHandler{
		service: composition.NewCompositionService(docService, cfg),
	}
}

func (h CompositionHandler) Create(c *gin.Context) {
	ehrId := c.Param("ehrid")
	if h.service.Doc.ValidateId(ehrId, types.EHR) == false {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panic("Cant close body request")
		}
	}(c.Request.Body)

	var request model.CompositionCreateRequest

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
	if errors.Is(err, errors.IsNotExist) {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	// Composition document creating
	doc, err := h.service.Create(userId, &request)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Composition creating error"})
		return
	}

	h.respondWithDocOrHeaders(ehrId, doc, c)
}

func (h *CompositionHandler) respondWithDocOrHeaders(ehrId string, doc *model.Composition, c *gin.Context) {
	uid := doc.Uid.Value
	c.Header("Location", h.service.Cfg.BaseUrl+"/v1/ehr/"+ehrId+"/composition/"+uid)
	c.Header("ETag", uid)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}
