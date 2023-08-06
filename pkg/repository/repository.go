package repository

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoList interface {
	CreateTask(task models.Task) (int, error)
	ReadTask(id int) (models.Task, error)
	UpdateStatus(id int) error
	Delete(id int) error
	UpdateTask(id int, task models.Task) error
}

type Repository struct {
	ToDoList
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		ToDoList: NewToDoListMongoDB(db),
	}
}
