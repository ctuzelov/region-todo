package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ctuzelov/region-todo/cmd/server"
	"github.com/ctuzelov/region-todo/pkg/handler"
	"github.com/ctuzelov/region-todo/pkg/repository"
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title TODO TASKS API
// @version 1.0
// @description This is a sample todo server.

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	// * начинаем выполнять инициализацию конфигурации сервера
	if err := server.InitConfig(); err != nil {
		log.Fatalf("error %s occured while initializating configs", err)
	}

	// ? хост и пароль нужны в случае если запускать с переменными окружения
	db, err := repository.NewMongoDB(repository.Config{
		Driver:   viper.GetString("db.driver"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
	})

	// * Закрываем соединение, когда функция вернется
	defer func() {
		if err = db.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to close the database connection:", err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()

	// * инициализация слоев и сервера
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	srv := new(server.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Fatalf("error %s in runnig the http server", err)
		}
	}()
	log.Printf("listening on http://localhost:" + viper.GetString("port"))

	// * Контролируемая остановка приложения: ожидание сигнала завершения,
	// * завершение сервера и обработка ошибок.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
