package repository

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoList interface {
	CreateTask(task models.Task) (int, error)
}

type Repository struct {
	ToDoList
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		ToDoList: NewToDoListMongoDB(db),
	}
}
