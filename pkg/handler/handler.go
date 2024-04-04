package handler

import (
	"github.com/spf13/viper"
	"net/http"

	"github.com/OTumanov/go_final_project/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/api")
	{
		auth.POST("/sign", h.login)
	}

	api := router.Group("/api", h.authMiddleware)
	{
		api.GET("/nextdate", h.nextDate)

		api.POST("/task", h.createTask)
		api.GET("/task", h.getTaskById)
		api.GET("/tasks", h.getTasks)
		api.PUT("/task", h.updateTask)
		api.POST("/task/done", h.taskDone)
		api.DELETE("/task", h.deleteTask)
	}

	static := router.Group("/")
	{
		router.GET("/", h.indexPage)

		static.StaticFS("./css", http.Dir(viper.Get("WEBDir").(string)+"/css"))
		static.StaticFS("./js", http.Dir(viper.Get("WEBDir").(string)+"/js"))
		router.StaticFile("/index.html", "./web/index.html")
		router.StaticFile("/login.html", "./web/login.html")
		router.StaticFile("/favicon.ico", "./web/favicon.ico")

	}
	return router
}

func (h *Handler) indexPage(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "./web/index.html")
}
