package api

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/processing"
	userModel "hms/gateway/pkg/user/model"
)

type DirectoryService interface {
	NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error)
	Create(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error
	Update(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error
	Delete(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) error
	GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error)
	GetByVersion(ctx context.Context, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) (*model.Directory, error)

	//GetFolder(ctx context.Context, d *model.Directory, path string) (*model.Directory, error)
	//Get(ctx context.Context, userID string, cID string) (*model.ContributionResponse, error)
	//Validate(ctx context.Context, c *model.Contribution, template helper.Searcher) (bool, error)
}

type DirectoryHandler struct {
	service            DirectoryService
	userService        UserService
	templateService    TemplateService
	compositionService CompositionService
	baseURL            string
}

func NewDirectoryHandler(cS DirectoryService, uS UserService, tS TemplateService, compS CompositionService, baseURL string) *DirectoryHandler {
	return &DirectoryHandler{
		service:            cS,
		userService:        uS,
		templateService:    tS,
		compositionService: compS,
		baseURL:            baseURL,
	}
}

// Create
// @Summary      Create CONTRIBUTION
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/ehr.html#tag/CONTRIBUTION/operation/contribution_create
// @Tags         CONTRIBUTION
// @Param        Authorization  header    string  true   "Bearer AccessToken"
// @Param        AuthUserId     header    string  true   "UserId UUID"
// @Param        Prefer         header    string     true  "Request header to indicate the preference over response details. The response will contain the entire resource when the Prefer header has a value of return=representation."  Enums: ("return=representation", "return=minimal") default("return=minimal")
// @Header       201            {string}  Etag   "The ETag (i.e. entity tag) response header is the contribution_uid identifier, enclosed by double quotes. Example: \"0826851c-c4c2-4d61-92b9-410fb8275ff0\""
// @Header       201            {string}  Location   "{baseUrl}/ehr/{ehr_id}/contribution/{contribution_uid}"
// @Header       201            {string}  RequestID  "Request identifier"
// @Accept       json
// @Produce      json
// @Success      201  {object}  model.ContributionResponse  "Is returned when the CONTRIBUTION was successfully created."
// @Failure      400  {object} model.ErrorResponse "Is returned when the request URL or body could not be parsed or has invalid content (e.g. invalid {ehr_id}, or either the body of the request not be converted to a valid CONTRIBUTION object, or the modification type doesn’t match the operation - i.e. first version of a composition with MODIFICATION)."
// @Failure      409  "Is returned when a resource with same identifier(s) already exists."
// @Failure      500  "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/contribution [post]
func (h *DirectoryHandler) Create(ctx *gin.Context)       {}
func (h *DirectoryHandler) Update(ctx *gin.Context)       {}
func (h *DirectoryHandler) Delete(ctx *gin.Context)       {}
func (h *DirectoryHandler) GetByTime(ctx *gin.Context)    {}
func (h *DirectoryHandler) GetByVersion(ctx *gin.Context) {}
