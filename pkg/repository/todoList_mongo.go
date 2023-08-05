package repository

import (
	"context"
	"log"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoListMongoDB struct {
	db *mongo.Client
}

func NewToDoListMongoDB(db *mongo.Client) *ToDoListMongoDB {
	return &ToDoListMongoDB{db}
}

func (r *ToDoListMongoDB) CreateTask(task models.Task) (int, error) {
	// Get a handle to the "users" collection.
	collection := r.db.Database("tasks").Collection("todo_list")

	// Insert the document into the "users" collection.
	inserted, err := collection.InsertOne(context.Background(), task)
	insertedID := inserted.InsertedID.(primitive.ObjectID)
	if err != nil {
		return 0, err
	}

	// Define a filter to match all documents (an empty filter means all documents).
	filter := bson.M{}

	// Execute the query and get the cursor to the results.
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor to access the documents.
	var tasks []models.Task
	for cursor.Next(context.Background()) {
		var todo models.Task
		if err := cursor.Decode(&todo); err != nil {
			return 0, err
		}
		tasks = append(tasks, todo)
	}

	// Handle any errors that occurred during iteration.
	if err := cursor.Err(); err != nil {
		log.Fatal("Cursor error:", err)
	}

	return util.ByteArrayToInt(insertedID), nil
}
