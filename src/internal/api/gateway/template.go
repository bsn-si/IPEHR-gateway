package gateway

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/template"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
)

type TemplateService interface {
	helper.Finder
	Parser(version model.ADLVer) (template.ADLParser, error)
	GetByID(ctx context.Context, userID, systemID, templateID string) (*model.Template, error)
	Store(ctx context.Context, userID, systemID, reqID string, m *model.Template) error
	GetList(ctx context.Context, userID, systemID string) ([]*model.Template, error)
}

type TemplateHandler struct {
	service TemplateService
	baseURL string
}

func NewTemplateHandler(templateService TemplateService, baseURL string) *TemplateHandler {
	return &TemplateHandler{
		service: templateService,
		baseURL: baseURL,
	}
}

// Get
//
//	@Summary		Get a template
//	@Description	Retrieves the ADL 1.4 operational template (OPT) identified by {template_id} identifier.
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list
//	@Tags			DEFINITION
//	@Produce		application/xml
//	@Produce		application/openehr.wt+json
//	@Param			template_id		path		string	false	"Template identifier. Example: Vital Signs"
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Param			EhrSystemId		header		string	false	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200				{string}	[]byte
//	@Failure		400				"Is returned when the request has invalid content."
//	@Failure		404				"Is returned when a stored query with {qualified_query_name} and {version} does not exist."
//	@Failure		406				"Is returned when template with certain ID created with other accept header"
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/template/adl1.4/{template_id} [get]
func (h *TemplateHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	tID := c.Param("template_id")

	t, err := h.service.GetByID(c, userID, systemID, tID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if c.GetHeader("Accept") != t.MimeType {
		c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": "Accept should be: " + t.MimeType})
		return
	}

	c.Data(http.StatusOK, t.MimeType, t.Body)
}

// Get
//
//	@Summary		Get a list of templates
//	@Description	List the available ADL 1.4 operational templates (OPT) on the system.
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_list
//	@Tags			DEFINITION
//	@Produce		json
//	@Param			Authorization	header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string	true	"UserId"
//	@Success		200				{object}	[]model.Template
//	@Failure		400				"Is returned because of invalid content."
//	@Failure		500				"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/template/adl1.4 [get]
func (h *TemplateHandler) ListStored(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	l, err := h.service.GetList(c, userID, systemID)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(l) == 0 {
		l = make([]*model.Template, 0)
	}

	c.JSON(http.StatusOK, l)
}

// Store
//
//	@Summary		Store a template
//	@Description	Upload a new ADL 1.4 operational template (OPT).
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/ADL1.4/operation/definition_template_adl1.4_upload
//	@Tags			DEFINITION
//	@Param			Authorization	header		string		true	"Bearer AccessToken"
//	@Param			AuthUserId		header		string		true	"UserId"
//	@Param			Prefer			header		string		true	"Request header to indicate the preference over response details. The response will contain the entire resource when the Prefer header has a value of return=representation."	Enums:	("return=representation", "return=minimal")	default("return=minimal")
//	@Header			201				{string}	Location	"{baseUrl}/definition/template/adl1.4/{template_id}"
//	@Header			201				{string}	RequestID	"Request identifier"
//	@Accept			application/xml
//	@Produce		text/plain
//	@Produce		application/xml
//	@Success		201	{object}	model.Template	"Is returned when the query was successfully uploaded."
//	@Failure		400	"Is returned when unable to upload a template, because of invalid content."
//	@Failure		409	"Is returned when a template with same {template_id} (at given version, if supplied) already exists."
//	@Failure		500	"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/template/adl1.4 [post]
func (h *TemplateHandler) Store(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}

	p, err := h.service.Parser(model.VerADL1_4)
	if err != nil {
		log.Printf("Template service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if !p.Validate(data, model.ADLTypeXML) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is not valid"})
		return
	}

	m, err := p.ParseWithFill(data, model.ADLTypeXML)
	if err != nil {
		log.Printf("Template service parse error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	reqID := c.GetString("reqID")
	if err := h.service.Store(c, userID, systemID, reqID, m); err != nil {
		log.Printf("Template service store error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", h.baseURL+"/definition/template/"+m.VerADL+"/"+url.QueryEscape(m.TemplateID))

	if c.Request.Header.Get("Prefer") == "return=representation" {
		c.Data(http.StatusCreated, "application/xml", data)
		return
	}

	c.Status(http.StatusCreated)
}
