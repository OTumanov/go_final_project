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
func (t *TaskService) CreateTask(task model.Task) (int64, error) {
	return t.repo.CreateTask(task)
}
func (t *TaskService) GetTasks(search string) (model.ListTasks, error) {
	return t.repo.GetTasks(search)
}
func (t *TaskService) GetTaskById(id string) (model.Task, error) {
	return t.repo.GetTaskById(id)
}
func (t *TaskService) UpdateTask(task model.Task) error {
	return t.repo.UpdateTask(task)
}
func (t *TaskService) DeleteTask(id string) error {
	return t.repo.DeleteTask(id)
}
func (t *TaskService) TaskDone(id string) error {
	return t.repo.TaskDone(id)
}
