package repository

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
}

func NewMongoDB(cnf Config) (*mongo.Client, error) {
	URI := fmt.Sprintf("%s://%s:%s", cnf.Driver, cnf.Username, cnf.Port)
	// * Установим параметры клиента
	clientOptions := options.Client().ApplyURI(URI)

	// * Подключимся к MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// * Проверка подключения
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %s", err.Error())
	}

	return client, nil
}
