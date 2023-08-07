package repository

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoTasks interface {
	CreateTask(task models.Task) (int, error)
	ReadTasks(status string) ([]models.Task, error)
	UpdateTaskStatus(id int) error
	DeleteTask(id int) error
	UpdateTask(id int, task models.Task) error
}

type Repository struct {
	ToDoTasks
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{
		ToDoTasks: NewToDoTasksMongoDB(db),
	}
}
