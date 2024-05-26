package repository

import (
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/jmoiron/sqlx"
)

// Слой репозитория
type RepositoryInterface interface{
	GetTasks() (task []structures.Task,  err error)
	GetTasksWithStr(str string) (task []structures.Task,  err error)
	GetTasksWithDate(date string) (task []structures.Task,  err error)
	CreateTask(task *structures.Task) (int64, error)
	DelTask(id int64) (err error)
	GetTask(id int64) (task structures.Task,  err error)
	UpdateTask(task *structures.Task) (err error)
}


type Repository struct {
	RepositoryInterface
}

// Конструктор репозитория
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{RepositoryInterface: NewTodoRepository(db)}
}