package main

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/controller"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/repo"
	"github.com/aksioto/awesome-task-exchange-system/cmd/accounting/usecase"
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
	//rabbitmqService.DeclareQueue(strings.Split(cfg.RabbitmqQueues, ",")) // need declare everywhere?

	// repo
	accountingRepo := repo.NewAccountingRepo(db)

	// usecase
	accountingUsecase := usecase.NewAccountingUsecase(accountingRepo)

	// controller
	accountingController := controller.NewAccountingController(accountingUsecase, rabbitmqService)

	// Async
	go rabbitmqService.Receive(accountingController.HandleNewTasks, "new_tasks")
	go rabbitmqService.Receive(accountingController.HandleShuffledTasks, "shuffled_tasks")

	r := gin.Default()
	// App middleware
	r.Use(middleware.NewAuthMiddleware())

	// Routes
	//r.POST("/create_new_task", accountingController.HandleCreateNewTask)

	tcpAddr := net.TCPAddr{Port: cfg.Port}
	log.Println("Server is starting on port:", cfg.Port)
	if err = http.ListenAndServe(tcpAddr.String(), r); err != nil {
		log.Fatalf("Failed to listen port: %o.\nError: %s", cfg.Port, err.Error())
	}
}
