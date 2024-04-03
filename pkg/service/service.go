package service

import (
	"github.com/OTumanov/go_final_project/pkg/model"
	"github.com/OTumanov/go_final_project/pkg/repository"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
	CreateTask(task model.Task) (int64, error)
	GetTasks(search string) (model.ListTasks, error)
	GetTaskById(id string) (model.Task, error)
	UpdateTask(task model.Task) error
	TaskDone(id string) error
	DeleteTask(id string) error
}

type Service struct {
	TodoTask
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoTask: NewTodoTaskService(repository.TodoTask),
	}
}
