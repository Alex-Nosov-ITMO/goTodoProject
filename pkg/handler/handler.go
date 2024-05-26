package handler

import (
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/service"
	"github.com/gin-gonic/gin"
)


// Слой обработчика
type HandlerInterface interface {
	InitRoutes() *gin.Engine
	GetTasks(*gin.Context) ()
	CreateTask(*gin.Context) ()
	GetTask(*gin.Context) ()
	DelTask(*gin.Context) ()
	UpdateTask(*gin.Context) ()
	NextDate(*gin.Context) ()
	DoneTask(*gin.Context) ()
	Login(*gin.Context) ()
}



type Handler struct {
	HandlerInterface
}


// Конструктор обработчика
func NewHandler(services *service.Service) *Handler {
	return &Handler{
		HandlerInterface: NewTodoHandler(services),
	}
}
