package service

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/repository"
)

type ToDoListService struct {
	repo repository.ToDoList
}

func NewToDoService(repo repository.ToDoList) *ToDoListService {
	return &ToDoListService{repo: repo}
}

func (s *ToDoListService) CreateTask(task models.Task) (int, error) {
	return s.repo.CreateTask(task)
}
