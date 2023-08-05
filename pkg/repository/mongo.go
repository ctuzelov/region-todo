package repository

import (
	"context"
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
	// Set client options.
	clientOptions := options.Client().ApplyURI("mongodb://root:123@localhost:2717")

	// Connect to MongoDB.
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %s", err.Error())
	}

	return client, nil
}
