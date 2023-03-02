package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

type ContributionService interface {
	NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error)
	GetByID(ctx context.Context, userID string, cID string) (*model.ContributionResponse, error)
	Store(ctx context.Context, req processing.RequestInterface, systemID string, user *userModel.UserInfo, c *model.Contribution) error
	Validate(ctx context.Context, c *model.Contribution, template helper.Searcher) (bool, error)
	Execute(ctx context.Context, req processing.RequestInterface, userID, ehrUUID string, c *model.Contribution, hComposition helper.Searcher) error
	PrepareResponse(ctx context.Context, systemID string, c *model.Contribution) (*model.ContributionResponse, error)
}

type ContributionHandler struct {
	service            ContributionService
	userService        UserService
	templateService    TemplateService
	compositionService CompositionService
	baseURL            string
}

func NewContributionHandler(cS ContributionService, uS UserService, tS TemplateService, compS CompositionService, baseURL string) *ContributionHandler {
	return &ContributionHandler{
		service:            cS,
		userService:        uS,
		templateService:    tS,
		compositionService: compS,
		baseURL:            baseURL,
	}
}

// Get
// Summary      Get CONTRIBUTION by id
// Description  Retrieves a CONTRIBUTION identified by {contribution_uid} and associated with the EHR identified by {ehr_id}.
// Description  https://specifications.openehr.org/releases/ITS-REST/latest/ehr.html#tag/CONTRIBUTION/operation/contribution_create
// Tags         CONTRIBUTION
// Produce      json
// Param        contribution_uid    path      string  false  "The CONTRIBUTION uid"
// Param        Authorization  header    string     true  "Bearer AccessToken"
// Param        AuthUserId     header    string     true  "UserId UUID"
// Param        EhrSystemId           header    string  true   "The identifier of the system, typically a reverse domain identifier"
// Success      200            {object}  model.ContributionResponse
// Failure      400            "Is returned when the request has invalid content."
// Failure      404            "Is returned when an EHR with {ehr_id}  does not exist, or when a CONTRIBUTION with {contribution_uid} does not exist"
// Failure      500            "Is returned when an unexpected error occurs while processing a request"
// Router       /ehr/{ehr_id}/contribution/{contribution_uid} [get]
func (h *ContributionHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// TODO do I need to checking userID by EHRid?
	cUID := c.Param("contribution_uid")
	if cUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contribution_uid is empty"})
		return
	}

	con, err := h.service.GetByID(c, userID, cUID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("Contribution service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, con)
}

// Create
// Summary      Create CONTRIBUTION
// Description  https://specifications.openehr.org/releases/ITS-REST/latest/ehr.html#tag/CONTRIBUTION/operation/contribution_create
// Tags         CONTRIBUTION
// Param        Authorization  header    string  true   "Bearer AccessToken"
// Param        AuthUserId     header    string  true   "UserId UUID"
// Param        Prefer         header    string     true  "Request header to indicate the preference over response details. The response will contain the entire resource when the Prefer header has a value of return=representation."  Enums: ("return=representation", "return=minimal") default("return=minimal")
// Header       201            {string}  Etag   "The ETag (i.e. entity tag) response header is the contribution_uid identifier, enclosed by double quotes. Example: \"0826851c-c4c2-4d61-92b9-410fb8275ff0\""
// Header       201            {string}  Location   "{baseUrl}/ehr/{ehr_id}/contribution/{contribution_uid}"
// Header       201            {string}  RequestID  "Request identifier"
// Accept       json
// Produce      json
// Success      201  {object} model.ContributionResponse "Is returned when the CONTRIBUTION was successfully created."
// Failure      400  {object} model.ErrorResponse "Is returned when the request URL or body could not be parsed or has invalid content (e.g. invalid {ehr_id}, or either the body of the request not be converted to a valid CONTRIBUTION object, or the modification type doesn’t match the operation - i.e. first version of a composition with MODIFICATION)."
// Failure      409  "Is returned when a resource with same identifier(s) already exists."
// Failure      500  "Is returned when an unexpected error occurs while processing a request"
// Router       /ehr/{ehr_id}/contribution [post]
func (h *ContributionHandler) Create(ctx *gin.Context) {
	errResponse := model.ErrorResponse{}

	userID := ctx.GetString("userID")
	if userID == "" {
		errResponse.SetMessage("Header required").
			AddError(errors.ErrFieldIsIncorrect("AuthUserId"))

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	ehrID := ctx.Param("ehrid")
	// TODO do we need checking userID by EHRid because it can be summary data with different users (GetEhrUUIDByUserID(c, userID, systemID))

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		errResponse.SetMessage("Incorrect param").
			AddError(errors.ErrFieldIsIncorrect("ehr_id"))

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	c := &model.Contribution{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(c); err != nil {
		errResponse.SetMessage("Request body parsing error").AddError(err)

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	existContr, err := h.service.GetByID(ctx, userID, c.UID.Value)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		log.Println("contributionService.GetByID error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if existContr != nil {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}

	systemID := ctx.GetString("ehrSystemID")

	templateSearcher := helper.NewSearcher(ctx, userID, systemID, ehrUUID.String())
	if ok, err := h.service.Validate(ctx, c, templateSearcher.UseService(h.templateService)); !ok {
		errResponse.SetMessage("Validation failed").AddError(err)

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	reqID := ctx.GetString("reqID")

	procRequest, err := h.service.NewProcRequest(reqID, userID, ehrUUID.String(), processing.RequestContributionCreate)
	if err != nil {
		log.Println("contributionService.Create error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	if err := h.service.Execute(ctx, procRequest, userID, ehrUUID.String(), c, templateSearcher.UseService(h.compositionService)); err != nil {
		errResponse.SetMessage("Execute failed").AddError(err)

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	userInfo, err := h.userService.Info(ctx, userID, systemID)
	if err != nil {
		log.Println("userService.Info error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.service.Store(ctx, procRequest, systemID, userInfo, c); err != nil {
		log.Printf("Contribution service store error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("Contribution procRequest commit error:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("Location", fmt.Sprintf("%s/%s/contribution/%s", h.baseURL, ehrUUID.String(), c.UID.Value))
	ctx.Header("ETag", fmt.Sprintf("\"%s\"", c.UID.Value))

	prefer := ctx.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		cR, err := h.service.PrepareResponse(ctx, systemID, c)
		if err != nil {
			log.Println("Contribution PrepareResponse error:", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, cR)
		return
	}

	ctx.Status(http.StatusCreated)
}
