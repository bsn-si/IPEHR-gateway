package stat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetPatientsCount(ctx context.Context, period string) (uint64, error)
	GetDocumentsCount(ctx context.Context, period string) (uint64, error)
}

// nolint
type StatHandler struct {
	service Service
}

func NewStatHandler(svc Service) *StatHandler {
	return &StatHandler{
		service: svc,
	}
}

type Stat struct {
	Patients  uint64 `json:"patients"`
	Documents uint64 `json:"documents"`
	Time      uint64 `json:"time"`
}

type ResponsePeriod struct {
	Type string `json:"type"`
	Data Stat   `json:"data"`
}

type ResponseTotal struct {
	Type  string `json:"type"`
	Data  Stat   `json:"data"`
	Month Stat   `json:"month"`
}

// GetStatPerMonth
// @Summary      Get IPEHR statistics per month
// @Description  Retrieve the IPEHR statistics per month
// @Tags         Stat
// @Produce      json
// @Param        period  path      string  false  "Month in YYYYYMM format. Example: 202201"
// @Success      200     {object}  ResponsePeriod
// @Failure      500     "Is returned when an unexpected error occurs while processing a request"
// @Router       /{period} [get]
func (h *StatHandler) GetStat(c *gin.Context) {
	period := c.Param("period")

	patientsCount, err := h.service.GetPatientsCount(c.Request.Context(), period)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, period)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	documentsCount, err := h.service.GetDocumentsCount(c.Request.Context(), period)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, period)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	periodInt, _ := strconv.Atoi(period)

	resp := ResponsePeriod{
		Type: "PERIOD",
		Data: Stat{
			Patients:  patientsCount,
			Documents: documentsCount,
			Time:      uint64(periodInt),
		},
	}

	c.JSON(http.StatusOK, resp)
}

// GetStat
// @Summary      Get IPEHR statistics total
// @Description  Retrieve the IPEHR statistics total and current month
// @Tags         Stat
// @Produce      json
// @Success      200     {object}  ResponseTotal
// @Failure      500     "Is returned when an unexpected error occurs while processing a request"
// @Router       / [get]
func (h *StatHandler) GetTotal(c *gin.Context) {
	currMonth := fmt.Sprintf("%d%02d", time.Now().Year(), time.Now().Month())

	patientsCurrMonth, err := h.service.GetPatientsCount(c.Request.Context(), currMonth)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, currMonth)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	documentsCurrMonth, err := h.service.GetDocumentsCount(c.Request.Context(), currMonth)
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, currMonth)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	patientsTotal, err := h.service.GetPatientsCount(c.Request.Context(), "")
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, currMonth)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	documentsTotal, err := h.service.GetDocumentsCount(c.Request.Context(), "")
	if err != nil {
		log.Printf("service.GetPatientsCount error: %v period: %s", err, currMonth)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	currMonthInt, _ := strconv.Atoi(currMonth)

	resp := ResponseTotal{
		Type: "LATEST",
		Data: Stat{
			Patients:  patientsTotal,
			Documents: documentsTotal,
			Time:      uint64(time.Now().Unix()),
		},
		Month: Stat{
			Patients:  patientsCurrMonth,
			Documents: documentsCurrMonth,
			Time:      uint64(currMonthInt),
		},
	}

	c.JSON(http.StatusOK, resp)
}
