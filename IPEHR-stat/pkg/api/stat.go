package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ipehr/stat/pkg/localDB"
	"ipehr/stat/pkg/service/stat"
)

type StatHandler struct {
	service *stat.Service
}

func NewStatHandler(db *localDB.DB) *StatHandler {
	return &StatHandler{
		service: stat.NewService(db),
	}
}

func (h *StatHandler) PatientsCount(c *gin.Context) {
	period := c.Param("period")

	count, err := h.service.GetPatientsCount(period)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, period)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "%d", count)
}

func (h *StatHandler) DocumentsCount(c *gin.Context) {
	period := c.Param("period")

	count, err := h.service.GetDocumentsCount(period)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, period)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.String(http.StatusOK, "%d", count)
}
