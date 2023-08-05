package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ctuzelov/region-todo/internal/repository"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Name string `bson:"name"`
}

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

	// Get a handle to the "users" collection.
	collection := db.Database("user").Collection("users")

	user := User{
		Name: "John Doe",
	}

	// Insert the document into the "users" collection.
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal("Error inserting document:", err)
	}

	// Define a filter to match all documents (an empty filter means all documents).
	filter := bson.M{}

	// Execute the query and get the cursor to the results.
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Error executing the query:", err)
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor to access the documents.
	var users []User
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal("Error decoding document:", err)
		}
		users = append(users, user)
	}

	// Handle any errors that occurred during iteration.
	if err := cursor.Err(); err != nil {
		log.Fatal("Cursor error:", err)
	}

	// Print the fetched users.
	fmt.Println("Fetched users:")
	for _, user := range users {
		fmt.Println(user.Name)
	}
}

func initConfig() error {
	viper.AddConfigPath("pkg/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
