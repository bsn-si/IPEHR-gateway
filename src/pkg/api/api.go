package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/groupAccess"
	"hms/gateway/pkg/infrastructure"
)

// @title        IPEHR Gateway API
// @version      0.2
// @description  The IPEHR Gateway is an openEHR compliant EHR server implementation that stores encrypted medical data in a Filecoin distributed file storage.

// @contact.name   API Support
// @contact.url    https://bsn.si/blockchain
// @contact.email  support@bsn.si

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      gateway.ipehr.org
// host      localhost:8080
// @BasePath  /v1

type API struct {
	Ehr         *EhrHandler
	EhrStatus   *EhrStatusHandler
	Composition *CompositionHandler
	Query       *QueryHandler
	GroupAccess *GroupAccessHandler
	Request     *RequestHandler
	User        *UserHandler
}

func New(cfg *config.Config, infra *infrastructure.Infra) *API {
	docService := service.NewDefaultDocumentService(cfg, infra)
	groupAccessService := groupAccess.NewService(docService, cfg.DefaultGroupAccessID, cfg.DefaultUserID)

	return &API{
		Ehr:         NewEhrHandler(docService, cfg.BaseURL),
		EhrStatus:   NewEhrStatusHandler(docService, cfg.BaseURL),
		Composition: NewCompositionHandler(docService, groupAccessService, cfg.BaseURL),
		Query:       NewQueryHandler(docService),
		GroupAccess: NewGroupAccessHandler(docService, groupAccessService, cfg.BaseURL),
		Request:     NewRequestHandler(docService),
		User:        NewUserHandler(cfg, infra, docService.Proc),
	}
}

func (a *API) Build() *gin.Engine {
	r := gin.New()

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(404)
	})

	r.Use(requestID)

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[GIN] %19s | %6s | %3d | %13v | %15s | %-7s %#v %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Keys["reqId"],
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))

	v1 := r.Group("v1")
	ehr := v1.Group("ehr")
	access := v1.Group("access")
	query := v1.Group("query")
	requests := v1.Group("requests")
	user := v1.Group("user")

	a.setRedirections(r).
		buildUserAPI(user).
		buildEhrAPI(ehr).
		buildGroupAccessAPI(access).
		buildQueryAPI(query).
		buildRequestsAPI(requests)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (a *API) buildEhrAPI(r *gin.RouterGroup) *API {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(Recovery, app_errors.ErrHandler)
	r.Use(auth(a))
	r.Use(ehrSystemID)
	r.POST("", a.Ehr.Create)
	r.GET("", a.Ehr.GetBySubjectIDAndNamespace)
	r.PUT("/:ehrid", a.Ehr.CreateWithID)
	r.GET("/:ehrid", a.Ehr.GetByID)
	r.PUT("/:ehrid/ehr_status", a.EhrStatus.Update)
	r.GET("/:ehrid/ehr_status/:versionid", a.EhrStatus.GetByID)
	r.GET("/:ehrid/ehr_status", a.EhrStatus.GetStatusByTime)
	r.POST("/:ehrid/composition", a.Composition.Create)
	r.GET("/:ehrid/composition/:version_uid", a.Composition.GetByID)
	r.DELETE("/:ehrid/composition/:preceding_version_uid", a.Composition.Delete)
	r.PUT("/:ehrid/composition/:versioned_object_uid", a.Composition.Update)

	return a
}

func (a *API) buildGroupAccessAPI(r *gin.RouterGroup) *API {
	r.Use(auth(a))
	r.GET("/group/:group_id", a.GroupAccess.Get)
	r.POST("/group", a.GroupAccess.Create)

	return a
}

func (a *API) buildQueryAPI(r *gin.RouterGroup) *API {
	r.Use(auth(a))
	r.POST("/aql", a.Query.ExecPost)

	return a
}

func (a *API) buildRequestsAPI(r *gin.RouterGroup) *API {
	r.Use(auth(a))
	r.GET("/", a.Request.GetAll)
	r.GET("/:reqId", a.Request.GetByID)

	return a
}

func (a *API) buildUserAPI(r *gin.RouterGroup) *API {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(ehrSystemID)
	r.POST("/register", a.User.Register)
	r.POST("/login", a.User.Login)
	r.GET("/refresh", a.User.RefreshToken)

	r.Use(auth(a))
	r.POST("/logout", a.User.Logout)
	return a
}

func (a *API) setRedirections(r *gin.Engine) *API {
	redirect := func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "v1/")
	}

	r.GET("/", redirect)
	r.HEAD("/", redirect)

	return a
}
