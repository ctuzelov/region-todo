package service

import (
	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/repository"
)

type ToDoListService struct {
	repo repository.ToDoTasks
}

func NewToDoService(repo repository.ToDoTasks) *ToDoListService {
	return &ToDoListService{repo: repo}
}

func (s *ToDoListService) CreateTask(task models.Task) (int, error) {
	return s.repo.CreateTask(task)
}

func (s *ToDoListService) ReadTasks(status string) ([]models.Task, error) {
	return s.repo.ReadTasks(status)
}

func (s *ToDoListService) DeleteTask(id int) error {
	return s.repo.DeleteTask(id)
}

func (s *ToDoListService) UpdateTaskStatus(id int) error {
	return s.repo.UpdateTaskStatus(id)
}

func (s *ToDoListService) UpdateTask(id int, task models.Task) error {
	return s.repo.UpdateTask(id, task)
}
