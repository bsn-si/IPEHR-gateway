package api

import (
	"fmt"
	"hms/gateway/pkg/docs/service/storedQuery"
	"hms/gateway/pkg/docs/service/user"
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
	StoredQuery *StoredQueryHandler
	//GroupAccess *GroupAccessHandler
	DocAccess *DocAccessHandler
	Request   *RequestHandler
	User      *UserHandler
}

func New(cfg *config.Config, infra *infrastructure.Infra) *API {
	docService := service.NewDefaultDocumentService(cfg, infra)
	groupAccessService := groupAccess.NewService(docService, cfg.DefaultGroupAccessID, cfg.DefaultUserID)
	storedQueryService := storedQuery.NewService(docService)
	userService := user.NewUserService(infra, docService.Proc)

	return &API{
		Ehr:         NewEhrHandler(docService, cfg.BaseURL),
		EhrStatus:   NewEhrStatusHandler(docService, cfg.BaseURL),
		Composition: NewCompositionHandler(docService, groupAccessService, cfg.BaseURL),
		Query:       NewQueryHandler(docService),
		StoredQuery: NewStoredQueryHandler(storedQueryService),
		//GroupAccess: NewGroupAccessHandler(docService, groupAccessService, cfg.BaseURL),
		DocAccess: NewDocAccessHandler(docService),
		Request:   NewRequestHandler(docService),
		User:      NewUserHandler(userService),
	}
}

func (a *API) Build() *gin.Engine {
	return a.setupRouter(
		a.buildUserAPI(),
		a.buildEhrAPI(),
		a.buildAccessAPI(),
		//a.buildGroupAccessAPI(),
		a.buildQueryAPI(),
		a.buildStoredQueryAPI(),
		a.buildRequestsAPI(),
	)
}

type handlerBuilder func(r *gin.RouterGroup)

func (a *API) setupRouter(apiHandlers ...handlerBuilder) *gin.Engine {
	r := gin.New()

	setRedirections(r)

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(404)
	})

	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("[GIN] %19s | %6s | %3d | %13v | %15s | %-7s %#v %s\n",
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			param.Keys["reqID"],
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))

	r.Use(requestID)

	v1 := r.Group("v1")
	for _, b := range apiHandlers {
		b(v1)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (a *API) buildEhrAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("ehr")
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
	}
}
func (a *API) buildAccessAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("access")
		r.Use(auth(a))
		//r.GET("/group/:group_id", a.GroupAccess.Get)
		//r.POST("/group", a.GroupAccess.Create)

		r.POST("/document", a.DocAccess.Set)
		r.GET("/document/", a.DocAccess.List)
	}
}

func (a *API) buildQueryAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("query")
		r.Use(auth(a))
		r.POST("/aql", a.Query.ExecPost)
	}
}

func (a *API) buildStoredQueryAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("definition").Group("query")
		r.Use(auth(a))
		r.GET(":qualifiedQueryName", a.StoredQuery.Get)
		r.GET("/", a.StoredQuery.Get) // this need because GIN breaking API tests if qualifiedQueryName not set
	}
}

func (a *API) buildRequestsAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("requests")
		r.Use(auth(a, "userRegister"))
		r.GET("/:reqID", a.Request.GetByID)
		r.GET("/", a.Request.GetAll)
	}
}

func (a *API) buildUserAPI() handlerBuilder {
	return func(r *gin.RouterGroup) {
		r = r.Group("user")
		r.Use(gzip.Gzip(gzip.DefaultCompression))
		r.Use(ehrSystemID)
		r.POST("/register", a.User.Register)
		r.POST("/login", a.User.Login)
		r.GET("/refresh", a.User.RefreshToken)

		r.Use(auth(a))
		r.POST("/logout", a.User.Logout)
	}
}

func setRedirections(r *gin.Engine) *gin.Engine {
	redirect := func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "v1/")
	}

	r.GET("/", redirect)
	r.HEAD("/", redirect)

	return r
}
