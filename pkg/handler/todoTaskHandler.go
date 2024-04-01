package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/OTumanov/go_final_project/pkg/model"
)

func (h *Handler) NextDate(c *gin.Context) {
	var nd model.NextDate

	if c.ShouldBindQuery(&nd) == nil {
		logrus.Println(nd.Now)
		logrus.Println(nd.Date)
		logrus.Println(nd.Repeat)
	}

	fmt.Println(nd)

	str, err := h.service.TodoTask.NextDate(nd)
	logrus.Warn(str)
	if err != nil {
		c.Status(http.StatusBadRequest)
		logrus.Println(err)
	}
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(str))
}

func (h *Handler) createTask(c *gin.Context) {
}

func (h *Handler) deleteTask(c *gin.Context) {

}

func (h *Handler) updateTask(c *gin.Context) {

}

func (h *Handler) getTaskById(c *gin.Context) {

}

func (h *Handler) getTasks(c *gin.Context) {

}
