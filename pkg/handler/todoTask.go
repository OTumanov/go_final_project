package handler

import (
	"fmt"
	"net/http"

	"github.com/OTumanov/go_final_project/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) nextDate(c *gin.Context) {
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

func (h *Handler) getTaskById(c *gin.Context) {
	id := c.Query("id")
	task, err := h.service.TodoTask.GetTaskById(id)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, task)
}
func (h *Handler) getTasks(c *gin.Context) {
	search := c.Query("search")
	list, err := h.service.TodoTask.GetTasks(search)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, list)
}

func (h *Handler) updateTask(c *gin.Context) {
	var task model.Task

	if c.ShouldBindJSON(&task) == nil {
		logrus.Println(fmt.Sprintf("Получили на обновление объект task со следующими данными: id: %s, date: %s, title: %s, comment: %s, repeat: %s", task.ID, task.Date, task.Title, task.Comment, task.Repeat))
	}
	_, err := h.service.TodoTask.GetTaskById(task.ID)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.TodoTask.UpdateTask(task)
	if err != nil {
		logrus.Error(err)
		NewResponseError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(200, gin.H{})
}
