package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/nextDate"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/service"
	"github.com/gin-gonic/gin"
)


type TodoHandler struct {
	srv service.ServiceInterface
}

// Конструктор связки сервиса и обработчика
func NewTodoHandler(services service.ServiceInterface) *TodoHandler {
	return &TodoHandler{
		srv: services,
	}
}

// Хендлеры
func (h *TodoHandler) GetTasks(c *gin.Context) {
	var tasks []structures.Task
	tasks, err := h.srv.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TodoHandler) CreateTask(c *gin.Context) {
	var task structures.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	id, err := h.srv.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *TodoHandler) GetTask(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}
	validId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id: %s, error: %v", id, err)})
		return
	}
	task, err := h.srv.GetTask(int64(validId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TodoHandler) DelTask(c *gin.Context) {

	id := c.Query("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id: %s, error: %v", id, err)})
		return
	}

	err = h.srv.DelTask(int64(validId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *TodoHandler) UpdateTask(c *gin.Context) {
	var task structures.Task
	
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	err := h.srv.UpdateTask(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{})
}	


func (h *TodoHandler) DoneTask(c *gin.Context) {
	id := c.Query("id")
	validId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid id: %s, error: %v", id, err)})
		return
	}


	var task structures.Task
	task, err = h.srv.GetTask(int64(validId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	switch task.Repeat {

	case "":
		err := h.srv.DelTask(int64(validId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	
	default:
		task.Date, err = nextDate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = h.srv.UpdateTask(&task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{})
		return
	}
}


func (h *TodoHandler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.Use(gin.Logger()) 
	router.Use(gin.Recovery())
	gin.SetMode(gin.ReleaseMode)
	
	
	router.StaticFS("/web", gin.Dir("./web", true))
    router.StaticFile("/favicon.ico", "./web/favicon.ico")
    router.StaticFile("/index.html", "./web/index.html")
    router.Static("/css", "./web/css")
    router.Static("/js", "./web/js")
    router.StaticFile("/login.html", "./web/login.html")

	router.GET("/", Index)



	api := router.Group("/api")
	{
		api.POST("/task", h.CreateTask)
		api.GET("/nextdate", h.NextDate)
		api.GET("/tasks", h.GetTasks)
		api.GET("/task", h.GetTask)
		api.PUT("/task", h.UpdateTask)
		api.DELETE("/task", h.DelTask)
		api.POST("/task/done", h.DoneTask)
	}


	return router
}


func Index(c *gin.Context) {
    c.File("./web/index.html")
}


func (h *TodoHandler) NextDate(c *gin.Context) {
	now := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	nowTime, err := time.Parse("20060102", now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	newDate, err := nextDate.NextDate(nowTime, date, repeat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.String(http.StatusOK, newDate)
}

