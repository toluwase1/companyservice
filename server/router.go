package server

import (
	"company-service/middleware"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func (s *Server) defineRoutes(router *gin.Engine) {
	apirouter := router.Group("/api/v1")
	apirouter.Use(middleware.Authenticate())
	apirouter.Use(middleware.IsCompanyOwner())
	apirouter.POST("/company", s.HandleCreateCompany())
	apirouter.PATCH("/company/:companyID", s.HandleUpdateCompany())
	apirouter.DELETE("/company/:companyID", s.handleDeleteCompany())
	apirouter.GET("/company/:companyID", s.HandleGetCompanyDetails())
}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()
	//staticFiles := "server/templates/static"
	//htmlFiles := "server/templates/*.html"
	//if s.Config.Env == "test" {
	//	_, b, _, _ := runtime.Caller(0)
	//	basepath := filepath.Dir(b)
	//	staticFiles = basepath + "/templates/static"
	//	htmlFiles = basepath + "/templates/*.html"
	//}
	//r.StaticFS("static", http.Dir(staticFiles))
	//r.LoadHTMLGlob(htmlFiles)

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	// setup cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.defineRoutes(r)

	return r
}
