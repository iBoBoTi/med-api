package server

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func (s *Server) defineRoutes(router *gin.Engine) {
	apirouter := router.Group("/api/v1")
	apirouter.POST("/auth/signup", s.HandleSignup())
	apirouter.POST("/auth/login", s.handleLogin())

	apirouter.GET("/fb/auth", s.handleFBLogin())
	apirouter.GET("fb/callback", s.fbCallbackHandler())

	apirouter.GET("/verifyEmail/:token", s.HandleVerifyEmail())
	apirouter.POST("/password/forgot", s.SendEmailForPasswordReset())
	apirouter.POST("/password/reset/:token", s.ResetPassword())

	authorized := apirouter.Group("/")
	authorized.Use(s.Authorize())
	authorized.POST("/logout", s.handleLogout())
	authorized.GET("/users", s.handleGetUsers())
	authorized.PUT("/me/update", s.handleUpdateUserDetails())
	authorized.GET("/me", s.handleShowProfile())
	authorized.POST("/user/medications", s.handleCreateMedication())
	authorized.GET("/user/medications/:id", s.handleGetMedDetail())
	authorized.GET("/user/medications", s.handleGetAllMedications())
	authorized.GET("/user/medications/next", s.handleGetNextMedication())

}

func (s *Server) setupRouter() *gin.Engine {
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "test" {
		r := gin.New()
		s.defineRoutes(r)
		return r
	}

	r := gin.New()
	r.StaticFS("static", http.Dir("server/templates/static"))
	r.LoadHTMLGlob("server/templates/*.html")

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
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.defineRoutes(r)

	return r
}
