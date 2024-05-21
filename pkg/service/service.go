package service

import (
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/repository"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
)

// Слой сервиса
type ServiceInterface interface{
	GetTasks() (task []structures.Task,  err error)
	CreateTask(task *structures.Task) (int64, error)
	DelTask(id int64) (err error)
	GetTask(id int64) (task structures.Task,  err error)
	UpdateTask(task *structures.Task) (err error)
}


type Service struct {
	ServiceInterface
}

// Конструктор сервиса
func NewService(repos *repository.Repository) *Service {
	return &Service{ServiceInterface: NewTodoService(repos)}
}