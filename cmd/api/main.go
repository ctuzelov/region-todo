package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ctuzelov/region-todo/cmd/server"
	"github.com/ctuzelov/region-todo/pkg/handler"
	"github.com/ctuzelov/region-todo/pkg/repository"
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error %s occured while initializating configs", err)
	}

	db, err := repository.NewMongoDB(repository.Config{
		Driver:   viper.GetString("db.driver"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
	})

	// Close the connection when the function returns.
	defer func() {
		if err = db.Disconnect(context.Background()); err != nil {
			log.Fatal("Failed to close the database connection:", err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}()

	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	srv := new(server.Server)
	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatalf("error %s in runnig the http server", err)
	}
	log.Printf("listening on http://localhost" + viper.GetString("port"))
}

func initConfig() error {
	viper.AddConfigPath("pkg/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
