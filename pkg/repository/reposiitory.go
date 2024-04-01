package repository

import (
	"github.com/OTumanov/go_final_project/pkg/model"
	"github.com/jmoiron/sqlx"
)

type TodoTask interface {
	NextDate(nd model.NextDate) (string, error)
}

type Repository struct {
	TodoTask
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoTask: NewTodoTaskSqlite(db),
	}
}