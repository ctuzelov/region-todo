package service

import "github.com/ctuzelov/region-todo/pkg/repository"

type Authorization interface {
}

type ToDoList interface {
}

type Service struct {
	Authorization
	ToDoList
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		ToDoList:      NewToDoService(repo.ToDoList),
	}
}
