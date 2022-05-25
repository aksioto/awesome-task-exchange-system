package main

import (
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/controller"
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/repo"
	"github.com/aksioto/awesome-task-exchange-system/cmd/auth/usecase"
	"github.com/aksioto/awesome-task-exchange-system/internal/helper"
	"github.com/aksioto/awesome-task-exchange-system/internal/service/rabbitmq"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
	"net/http"
)

type Credentials struct {
	ClientSecret string `env:"CLIENT_SECRET,required"`
}

type Config struct {
	Port                     int    `env:"PORT,required"`
	DbConnectionString       string `env:"DB_CONNECTION_STRING,required"`
	RabbitmqConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
	RabbitmqQueues           string `env:"RABBITMQ_QUEUES,required"`
	*Credentials
}

func main() {
	cfg := &Config{
		Credentials: &Credentials{},
	}

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
	rabbitmqService.DeclareExchange("user_stream")

	// repo
	authRepo := repo.NewAuthRepo(db)

	// usecase
	authUsecase := usecase.NewAuthUsecase(cfg.ClientSecret, authRepo)

	// controller
	authController := controller.NewAuthController(authUsecase, rabbitmqService)

	r := gin.Default()
	r.StaticFS("/app", http.Dir("public"))

	r.POST("/signin", authController.HandleSignIn)
	r.POST("/signup", authController.HandleSignUp)
	r.GET("/token", authController.HandleToken)

	tcpAddr := net.TCPAddr{Port: cfg.Port}
	log.Println("Server is starting on port:", cfg.Port)
	if err = http.ListenAndServe(tcpAddr.String(), r); err != nil {
		log.Fatalf("Failed to listen port: %o.\nError: %s", cfg.Port, err.Error())
	}
}
