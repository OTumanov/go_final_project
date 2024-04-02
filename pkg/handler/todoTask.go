package handler

import (
	"fmt"
	"net/http"

	"github.com/OTumanov/go_final_project/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) NextDate(c *gin.Context) {
	var nd model.NextDate

	if c.ShouldBindQuery(&nd) == nil {
		logrus.Println(fmt.Sprintf("Получили объект NextDate со следующими данными: date: %s, now: %s, repeat: %s", nd.Date, nd.Now, nd.Repeat))
	}
	str, err := h.service.TodoTask.NextDate(nd)
	if err != nil {
		logrus.Error(err)
		c.Status(http.StatusBadRequest)
	}
	c.Writer.WriteHeader(200)
	c.Writer.Write([]byte(str))
}

func (h *Handler) createTask(c *gin.Context) {
	var task model.Task
	if c.ShouldBindJSON(&task) == nil {
		logrus.Println(fmt.Sprintf("Получили объект task со следующими данными: date: %s, title: %s, comment: %s, repeat: %s", task.Date, task.Title, task.Comment, task.Repeat))
	}

	id, err := h.service.TodoTask.CreateTask(task)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{"id": id})

}

//func (h *Handler) deleteTask(c *gin.Context) {
//}
//
//func (h *Handler) updateTask(c *gin.Context) {
//}
//
//func (h *Handler) getTaskById(c *gin.Context) {
//}
//
//func (h *Handler) getTasks(c *gin.Context) {
//}
