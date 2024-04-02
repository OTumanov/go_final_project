package handler

import (
	"github.com/OTumanov/go_final_project/config"
	"github.com/OTumanov/go_final_project/pkg/service"

	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	static := router.Group("/")
	{
		static.StaticFS("./css", http.Dir(config.WEBDir+"/css"))
		static.StaticFS("./js", http.Dir(config.WEBDir+"/js"))
		router.StaticFS("/index.html", http.Dir(config.WEBDir))
		router.GET("./login.html", h.loginPage)
		router.GET("./favicon.ico", h.favicon)
	}

	api := router.Group("/api")
	{
		api.POST("/tasks", h.createTask)
		api.GET("/nextdate", h.NextDate)
	}
	return router
}

func (h *Handler) favicon(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "./web/favicon.ico")
}

func (h *Handler) loginPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "./web/login.html")
}
