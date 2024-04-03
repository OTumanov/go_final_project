package service

import (
	"github.com/OTumanov/go_final_project/pkg/model"
	"github.com/OTumanov/go_final_project/pkg/repository"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
	CreateTask(task model.Task) (int64, error)
	GetTasks(search string) (model.ListTodoTask, error)
	GetTaskById(id string) (model.Task, error)
	UpdateTask(task model.Task) error
}

type Service struct {
	TodoTask
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoTask: NewTodoTaskService(repository.TodoTask),
	}
}
