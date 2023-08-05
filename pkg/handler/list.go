package handler

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *Handler) createList(g *gin.Context) {
	// Get a handle to the "users" collection.
	collection := h.service.ToDoList.Database("user").Collection("users")

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

}

func (h *Handler) getListByID(g *gin.Context) {

}

func (h *Handler) updateListByID(g *gin.Context) {

}

func (h *Handler) deleteListByID(g *gin.Context) {

}
