package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/logger"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/db"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/server"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/handler"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/repository"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/service"
	_ "github.com/Alex-Nosov-ITMO/go_project_final/docs"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}


// @title TodoList API
// @version 1.0
// @description API Server for TodoList Application
// @host localhost:7540
// @BasePath /
func main() {


	// TODO: init logger
	logger.LogFile = logger.LoggerInit()
	defer logger.LogFile.Close()
	log.SetOutput(logger.LogFile)

	// TODO: connect to db
	dbConn, err := db.NewDbConnection()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Main: db connected")
	defer dbConn.Close()

	
	// TODO: init repo, service and handler
	repos := repository.NewRepository(dbConn.Db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// TODO: start server
	srv := new(server.Server)
	if err := srv.Run(os.Getenv("TODO_PORT"), handlers.InitRoutes()); err != nil {
		log.Fatalf("Main: failed to run server: %v", err)
	}
}
