package service

import (
	"github.com/OTumanov/go_final_project/pkg/model"
	"github.com/OTumanov/go_final_project/pkg/repository"
)

type TaskService struct {
	repo repository.TodoTask
}

func NewTodoTaskService(repo repository.TodoTask) *TaskService {
	return &TaskService{repo: repo}
}
func (t *TaskService) NextDate(nd model.NextDate) (string, error) {
	return t.repo.NextDate(nd)
}
