package api

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/service"
)

var AppConfig *config.Config

type API struct {
	Ehr       *EhrHandler
	EhrStatus *EhrStatusHandler

	fs http.FileSystem
}

func NewAPI() *API {
	docService := service.NewDefaultDocumentService()
	// Get config from the root of the project
	AppConfig = config.GetConfig("../../../config.json")

	return &API{
		Ehr:       NewEhrHandler(docService),
		EhrStatus: NewEhrStatusHandler(docService),
	}
}

func (a *API) Build() *gin.Engine {
	r := gin.New()

	r.NoRoute(func(c *gin.Context) {
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p := strings.TrimPrefix(r.URL.Path, "/v1"); len(p) < len(r.URL.Path) {
				if p == "/" || p == "" {
					c.Header("Cache-Control", "no-store, max-age=0")
				}
				c.FileFromFS(p, a.fs)
			} else {
				http.NotFound(w, r)
			}
		}).ServeHTTP(c.Writer, c.Request)
	})

	v1 := r.Group("v1")
	ehr := v1.Group("ehr")

	a.setRedirections(r).
		buildEhrAPI(ehr)

	return r
}

func (a *API) buildEhrAPI(r *gin.RouterGroup) *API {
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	//r.Use(Recovery, app_errors.ErrHandler)
	r.GET("/:ehrid", a.Ehr.GetById)
	//r.GET("/:ehrid", a.EHR.GetBySubjectIdAndNamespace)	//TODO

	// Other methods should be authorized
	r.Use(a.Auth)
	r.POST("/", a.Ehr.Create)
	r.PUT("/:ehrid", a.Ehr.CreateWithId)
	r.PUT("/:ehrid/ehr_status", a.EhrStatus.Update)
	r.GET("/:ehrid/ehr_status/:versionid", a.EhrStatus.GetById)
	r.GET("/v1/ehr/:ehrid/ehr_status", a.EhrStatus.GetStatus)

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
