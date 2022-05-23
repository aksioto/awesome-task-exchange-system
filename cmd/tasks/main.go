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
	"strings"
)

type Config struct {
	Port                     int    `env:"PORT,required"`
	DbConnectionString       string `env:"DB_CONNECTION_STRING,required"`
	RabbitmqConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
	RabbitmqQueues           string `env:"RABBITMQ_QUEUES,required"`
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
	rabbitmqService.DeclareExchanges(strings.Split(cfg.RabbitmqQueues, ","))

	// repo
	tasksRepo := repo.NewTasksRepo(db)

	// usecase
	tasksUsecase := usecase.NewTasksUsecase(tasksRepo)

	// controller
	tasksController := controller.NewTasksController(tasksUsecase, rabbitmqService)

	// Async
	go rabbitmqService.Receive(tasksController.HandleUserStream, "user_stream")

	r := gin.Default()
	//authorized := r.Group("/")
	// App middleware
	r.Use(middleware.NewAuthMiddleware())
	// Routes
	v1 := r.Group("v1")
	v1.POST("/add_new_task", tasksController.HandleAddNewTask)
	v1.POST("/shuffle_tasks", tasksController.HandleShuffleTasks)
	v1.POST("/complete_task", tasksController.HandleCompleteTask)
	v1.GET("/dashboard", tasksController.HandleDashboard)

	v2 := r.Group("v2")
	v2.POST("/add_new_task", tasksController.HandleAddNewTask)

	tcpAddr := net.TCPAddr{Port: cfg.Port}
	log.Println("Server is starting on port:", cfg.Port)
	if err = http.ListenAndServe(tcpAddr.String(), r); err != nil {
		log.Fatalf("Failed to listen port: %o.\nError: %s", cfg.Port, err.Error())
	}
}
