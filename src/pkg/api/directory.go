package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/helper"
	userModel "hms/gateway/pkg/user/model"
)

type DirectoryService interface {
	NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error)
	Create(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error
	Update(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error
	Delete(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) error
	GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error)
	GetByVersion(ctx context.Context, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) (*model.Directory, error)
	GetByID(ctx context.Context, userID string, versionUID string) (*model.Directory, error)
}

type DirectoryHandler struct {
	service     DirectoryService
	userService UserService
	indexer     Indexer
	baseURL     string
}

func NewDirectoryHandler(cS DirectoryService, uS UserService, indexer Indexer, baseURL string) *DirectoryHandler {
	return &DirectoryHandler{
		service:     cS,
		userService: uS,
		indexer:     indexer,
		baseURL:     baseURL,
	}
}

// Create
// @Summary      Create DIRECTORY
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/ehr.html#tag/DIRECTORY/operation/directory_create
// @Tags         DIRECTORY
// @Param        Authorization  header    string  true   "Bearer AccessToken"
// @Param        AuthUserId     header    string  true   "UserId UUID"
// @Param        Prefer         header    string     true  "Request header to indicate the preference over response details. The response will contain the entire resource when the Prefer header has a value of return=representation."  Enums: ("return=representation", "return=minimal") default("return=minimal")
// @Header       201            {string}  Etag   "The ETag (i.e. entity tag) response header is the version_uid identifier, enclosed by double quotes. Example: \"8849182c-82ad-4088-a07f-48ead4180515::openEHRSys.example.com::1\""
// @Header       201            {string}  Location   "{baseUrl}/ehr/{ehr_id}/directory/{version_uid}"
// @Header       201            {string}  RequestID  "Request identifier"
// @Accept       json
// @Produce      json
// @Success      201  {object}  model.Directory  "Is returned when the DIRECTORY was successfully created."
// @Failure      400  "Is returned when the request has invalid content"
// @Failure      404  "Is returned when an EHR with {ehr_id}  does not exist"
// @Failure      409  "Is returned when a resource with same identifier(s) already exists"
// @Failure      500  "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/directory [post]
func (h *DirectoryHandler) Create(ctx *gin.Context) {
	errResponse := model.ErrorResponse{}

	userID := ctx.GetString("userID")
	if userID == "" {
		errResponse.SetMessage("Header required").
			AddError(errors.ErrFieldIsIncorrect("AuthUserId"))

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	ehrID := ctx.Param("ehrid")

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		errResponse.SetMessage("Incorrect param").
			AddError(errors.ErrFieldIsIncorrect("ehr_id"))

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	systemID := ctx.GetString("ehrSystemID")
	searcher := helper.NewSearcher(ctx, userID, systemID, ehrUUID.String()).UseService(h.indexer)

	if !searcher.IsEhrBelongsToUser() {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	d := &model.Directory{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(d); err != nil {
		errResponse.SetMessage("Request body parsing error").AddError(err)

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	exist, err := h.service.GetByID(ctx, userID, d.UID.Value)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		log.Println("directoryService.GetByID error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if exist != nil {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}

	reqID := ctx.GetString("reqID")

	procRequest, err := h.service.NewProcRequest(reqID, userID, ehrID, processing.RequestDirectoryCreate)
	if err != nil {
		log.Println("directoryService.NewProcRequest error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
	}

	userInfo, err := h.userService.Info(ctx, userID)
	if err != nil {
		log.Println("userService.Info error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := h.service.Create(ctx, procRequest, systemID, &ehrUUID, userInfo, d); err != nil {
		errResponse.SetMessage("Execute failed").AddError(err)

		ctx.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if err := procRequest.Commit(); err != nil {
		log.Println("Directory procRequest commit error:", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Header("Location", fmt.Sprintf("%s/%s/directory/%s", h.baseURL, ehrUUID.String(), d.UID.Value))
	ctx.Header("ETag", fmt.Sprintf("\"%s\"", d.UID.Value))

	prefer := ctx.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		ctx.JSON(http.StatusCreated, d)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *DirectoryHandler) Update(ctx *gin.Context)       {}
func (h *DirectoryHandler) Delete(ctx *gin.Context)       {}
func (h *DirectoryHandler) GetByTime(ctx *gin.Context)    {}
func (h *DirectoryHandler) GetByVersion(ctx *gin.Context) {}
