package api

import (
	"context"
	"hms/gateway/pkg/errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/template"
)

type TemplateService interface {
	Parser(version model.VerADL) (template.ADLParser, error)
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
// @Tags         TEMPLATE
// @Produce      xml
// @Produce      application/openehr.wt+json
// @Param        template_id       path      string  false  "Template identifier. Example: Vital Signs"
// @Param        Authorization     header    string  true   "Bearer AccessToken"
// @Param        AuthUserId        header    string  true   "UserId UUID"
// @Success      200               {string}  []byte
// @Failure      400                               "Is returned when the request has invalid content."
// @Failure      404                               "Is returned when a stored query with {qualified_query_name} and {version} does not exist."
// @Failure      500                         "Is returned when an unexpected error occurs while processing a request"
// @Router       /definition/template/adl1.4/{template_id} [get]
func (h *TemplateHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	tID := c.Param("template_id") // TODO should have its own structure and validation method

	t, err := h.service.GetByID(c, userID, tID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	p, err := h.service.Parser(t.VerADL)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	needType := c.GetHeader("Accept")

	pType, err := p.AllowedType(needType)
	if err != nil {
		c.AbortWithStatus(http.StatusNotAcceptable)
		return
	}

	r, err := p.Parse(t.Content, pType)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, pType, r)
}
