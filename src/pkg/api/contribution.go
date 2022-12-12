package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
)

type ContributionService interface {
	GetByID(ctx context.Context, userID string, cID string) (*model.Contribution, error)
	Store(ctx context.Context) error
	Validate(ctx context.Context, data string) (bool, error)
	Execute(ctx context.Context) error
	//Rollback(ctx context.Context) error // TODO remove
}

type ContributionHandler struct {
	service     ContributionService
	userService UserService
	baseURL     string
}

func NewContributionHandler(cS ContributionService, uS UserService, baseURL string) *ContributionHandler {
	return &ContributionHandler{
		service:     cS,
		userService: uS,
		baseURL:     baseURL,
	}
}

// Get
// @Summary      Get CONTRIBUTION by id
// @Description  Retrieves a CONTRIBUTION identified by {contribution_uid} and associated with the EHR identified by {ehr_id}.
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/ehr.html#tag/CONTRIBUTION/operation/contribution_create
// @Tags         CONTRIBUTION
// @Produce      json
// @Param        contribution_uid    path      string  false  "The CONTRIBUTION uid"
// @Param        Authorization  header    string     true  "Bearer AccessToken"
// @Param        AuthUserId     header    string     true  "UserId UUID"
// @Param        EhrSystemId           header    string  true   "The identifier of the system, typically a reverse domain identifier"
// @Success      200            {string}  model.contribution
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      404            "Is returned when an EHR with {ehr_id}  does not exist, or when a CONTRIBUTION with {contribution_uid} does not exist"
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/contribution/{contribution_uid} [get]
func (h *ContributionHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header EhrSystemId is empty"})
		return
	}

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
// @Failure      400  {object} model.ContributionResponseWithError "Is returned when the request URL or body could not be parsed or has invalid content (e.g. invalid {ehr_id}, or either the body of the request not be converted to a valid CONTRIBUTION object, or the modification type doesn’t match the operation - i.e. first version of a composition with MODIFICATION)."
// @Failure      409  "Is returned when a resource with same identifier(s) already exists."
// @Failure      500  "Is returned when an unexpected error occurs while processing a request"
// @Router       /ehr/{ehr_id}/contribution [post]
func (h *ContributionHandler) Create(ctx *gin.Context) {
	userID := ctx.GetString("userID")
	if userID == "" {
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := ctx.GetString("ehrSystemID")
	if systemID == "" {
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "header EhrSystemId is empty"})
		return
	}

	data, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer ctx.Request.Body.Close()

	c := &model.Contribution{}
	if err := json.NewDecoder(ctx.Request.Body).Decode(c); err != nil {
		//log.Println("Contribution Create request unmarshal error", err)
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request body parsing error"})
		return
	}

	if ok, err := h.service.Validate(ctx, c); !ok {
		// TODO 409
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userInfo, err := h.userService.Info(ctx, userID)
	if err != nil {
		log.Println("userService.Info error: ", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	reqID := ctx.GetString("reqID")
	if err := h.service.Execute(ctx, userInfo, reqID, c); err != nil {
		// TODO StatusBadRequest shoul return corrent response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ""})
		return
	}

	if err := h.service.Store(ctx, userID, systemID, reqID, c); err != nil {
		log.Printf("Contribution service store error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// TODO UID into UUID?
	ctx.Header("Location", fmt.Sprintf("%s/%s/contribution/%s", h.baseURL, systemID, c.UID.Value))
	ctx.Header("ETag", fmt.Sprintf("\"%s\"", c.UID.Value))

	if ctx.Request.Header.Get("Prefer") == "return=representation" {
		ctx.Data(http.StatusCreated, "application/xml", data)
		return
	}

	ctx.Status(http.StatusCreated)
}
