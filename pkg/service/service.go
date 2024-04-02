package service

import (
	"github.com/OTumanov/go_final_project/pkg/model"
	"github.com/OTumanov/go_final_project/pkg/repository"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
}

type Service struct {
	TodoTask
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		TodoTask: NewTodoTaskService(repository.TodoTask),
	}
}
