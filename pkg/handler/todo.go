package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/nextDate"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/structures"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/middleware"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

	search := c.Query("search")

	tasks, err := h.srv.GetTasks(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TodoHandler) CreateTask(c *gin.Context) {
	var task structures.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		log.Printf("Handler: CreateTask: ShouldBindJSON: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid json")})
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
		log.Printf("Handler: GetTask: strconv.Atoi: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid id")})
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
		log.Printf("Handler: DelTask: strconv.Atoi: invalid id: %s, error: %s\n", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid id")})
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
		log.Printf("Handler: UpdateTask: ShouldBindJSON: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid json")})
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
		log.Printf("Handler: DoneTask: strconv.Atoi: invalid id: %s, error: %s\n", id, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid id")})
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
		api.GET("/nextdate", h.NextDate)
		api.POST("/signin", h.Login)

		api.Use(middleware.Auth)
		api.POST("/task", h.CreateTask)
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
		log.Println("Handler: NextDate: time.Parse: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("invalid date")})
	}

	newDate, err := nextDate.NextDate(nowTime, date, repeat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.String(http.StatusOK, newDate)
}

func (h *TodoHandler) Login(c *gin.Context) {

	var pass structures.Password
	if err := c.ShouldBindJSON(&pass); err != nil {
		log.Printf("Handler: Login: ShouldBindJSON: %s\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("invalid json")})
		return
	}

	realPassword := os.Getenv("TODO_PASSWORD")

	if pass.Password != realPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong password"})
		return
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"password": pass.Password,
	})

	signedToken, err := jwtToken.SignedString(structures.Secret)
	if err != nil {
		log.Printf("Handler: Login: jwtToken.SignedString: %s\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.New("failed to sign token")})
		return
	}
	log.Println("token: ", signedToken)
	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
