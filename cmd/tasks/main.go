package main

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/controller"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/repo"
	"github.com/aksioto/awesome-task-exchange-system/cmd/tasks/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	"github.com/aksioto/awesome-task-exchange-system/internal/middleware"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
	"net/http"
)

type Config struct {
	Port                     int    `env:"PORT,required"`
	DbConnectionString       string `env:"DB_CONNECTION_STRING,required"`
	RabbitmqConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
}

func main() {
	cfg := &Config{}

	err := helper.PrepareEnvConfig(cfg)
	if err != nil {
		log.Fatal("Error happened on IgniteConfig", err)
	}

	db, err := sqlx.Open("mysql", cfg.DbConnectionString)
	if err != nil {
		log.Fatalf("Connection failed. Error: %s", err.Error())
	}
	defer db.Close()

	// services
	rabbitmqService := rabbitmq.NewRabbitmqService(cfg.RabbitmqConnectionString)
	defer rabbitmqService.Close()

	// repo
	tasksRepo := repo.NewTasksRepo(db)

	// usecase
	tasksUsecase := usecase.NewTasksUsecase(tasksRepo)

	// controller
	tasksController := controller.NewTasksController(tasksUsecase, rabbitmqService)

	// Async
	go tasksController.StartReceiver()

	r := gin.Default()
	//authorized := r.Group("/")

	// App middleware
	r.Use(middleware.NewAuthMiddleware())

	// Routes
	r.POST("/create_new_task", tasksController.HandleCreateNewTask)
	r.POST("/shuffle_tasks", tasksController.HandleShuffleTasks)

	// For auth testing
	r.GET("/status", tasksController.HandleStatus)

	tcpAddr := net.TCPAddr{Port: cfg.Port}
	log.Println("Server is starting on port:", cfg.Port)
	if err = http.ListenAndServe(tcpAddr.String(), r); err != nil {
		log.Fatalf("Failed to listen port: %o.\nError: %s", cfg.Port, err.Error())
	}
}
