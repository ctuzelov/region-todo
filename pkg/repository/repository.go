package repository

import "go.mongodb.org/mongo-driver/mongo"

type Authorization interface {
}

type ToDoList interface {
	CreateList()
}

type Repository struct {
	Authorization
	ToDoList
}

func NewRepository(db *mongo.Client) *Repository {
	return &Repository{}
}
