package main

import (
	"log"
	"os"

	"github.com/Alex-Nosov-ITMO/go_project_final/internal/logger"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/db"
	"github.com/Alex-Nosov-ITMO/go_project_final/internal/server"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/handler"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/repository"
	"github.com/Alex-Nosov-ITMO/go_project_final/pkg/service"
	"github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
}

func main() {

	// оставил этот комментарий для будущей замены нахождения бд в корне проекта на миграции
	/*confDb := db.ConfigDB{
		User:     os.Getenv("TODO_USER"),
		Password: os.Getenv("TODO_PASSWORD"),
		Host:     os.Getenv("TODO_HOST"),
		Port:     os.Getenv("TODO_DB_PORT"),
		Dbname:   os.Getenv("TODO_DBNAME"),
	}

	if len(confDb.User) == 0 || len(confDb.Password) == 0 || len(confDb.Host) == 0 || len(confDb.Port) == 0 || len(confDb.Dbname) == 0 {
		log.Fatal("todo db params not set")
	}*/


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
