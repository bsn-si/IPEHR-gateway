package api

import (
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/storage"
)

// @title        IPEHR Gateway API
// @version      0.1
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
}

func New(cfg *config.Config) *API {
	sc := storage.NewConfig(cfg.StoragePath)
	storage.Init(sc)

	docService := service.NewDefaultDocumentService(cfg)

	return &API{
		Ehr:         NewEhrHandler(docService, cfg),
		EhrStatus:   NewEhrStatusHandler(docService, cfg),
		Composition: NewCompositionHandler(docService, cfg),
		Query:       NewQueryHandler(docService, cfg),
		GroupAccess: NewGroupAccessHandler(docService, cfg),
	}
}

func (a *API) Build() *gin.Engine {
	r := gin.New()

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(404)
	})

	v1 := r.Group("v1")
	ehr := v1.Group("ehr")
	access := v1.Group("access")
	query := v1.Group("query")

	a.setRedirections(r).
		buildEhrAPI(ehr).
		buildGroupAccessAPI(access).
		buildQueryAPI(query)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func (a *API) buildEhrAPI(r *gin.RouterGroup) *API {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(Recovery, app_errors.ErrHandler)

	// Other methods should be authorized
	r.Use(a.Auth)
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
	r.Use(a.Auth)
	r.GET("/group/:group_id", a.GroupAccess.Get)
	r.POST("/group", a.GroupAccess.Create)

	return a
}

func (a *API) buildQueryAPI(r *gin.RouterGroup) *API {
	r.Use(a.Auth)
	r.POST("/aql", a.Query.ExecPost)

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
