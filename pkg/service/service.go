package service

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/repository"
)

type ToDoList interface {
	CreateTask(task models.Task) (int, error)
}

type Service struct {
	ToDoList
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		NewToDoService(repo.ToDoList),
	}
}
