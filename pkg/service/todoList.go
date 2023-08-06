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

func (s *ToDoListService) ReadTask(id int) (models.Task, error) {
	return s.repo.ReadTask(id)
}

func (s *ToDoListService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ToDoListService) UpdateStatus(id int) error {
	return s.repo.UpdateStatus(id)
}

func (s *ToDoListService) UpdateTask(id int, task models.Task) error {
	return s.repo.UpdateTask(id, task)
}
