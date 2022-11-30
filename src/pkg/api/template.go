package api

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/template"
)

type TemplateService interface {
	Parser(version string) (template.ADLParser, error)
	GetByID(ctx context.Context, userID string, templateID string) (*model.Template, error)
}

type TemplateHandler struct {
	service TemplateService
	baseURL string
}

func NewTemplateHandler(templateService TemplateService, baseURL string) *TemplateHandler {
	return &TemplateHandler{
		//
		service: templateService,
		baseURL: baseURL,
	}
}

// Get
// @Summary      Get a template
// @Description  Retrieves the ADL 1.4 operational template (OPT) identified by {template_id} identifier.
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list
// @Tags         QUERY
// @Produce      xml
// @Produce      application/openehr.wt+json
// @Param        template_id       path      string  false  "Template identifier. Example: Vital Signs"
// @Param        Authorization     header    string  true   "Bearer AccessToken"
// @Param        AuthUserId        header    string  true   "UserId UUID"
// @Success      200               {object}  []model.Template
// @Failure      500                         "Is returned when an unexpected error occurs while processing a request"
// @Router       /definition/template/adl1.4/{template_id} [get]
func (h *TemplateHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	qName := c.Param("template_id")

	p, err := h.service.Parser(model.VerADL1_4)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	t, err := h.service.GetByID(c, userID, qName)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO convert to needed type

	c.JSON(http.StatusOK, t)
}
