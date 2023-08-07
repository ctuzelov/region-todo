package service

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type ToDoTasks interface {
	CreateTask(task models.Task) (int, error)
	ReadTasks(status string) ([]models.Task, error)
	UpdateTaskStatus(id int) error
	DeleteTask(id int) error
	UpdateTask(id int, task models.Task) error
}

type Service struct {
	ToDoTasks
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		NewToDoService(repo.ToDoTasks),
	}
}
