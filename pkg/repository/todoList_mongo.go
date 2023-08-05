package repository

import (
	"context"
	"log"

	"github.com/ctuzelov/region-todo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoListMongoDB struct {
	db *mongo.Client
}

func NewToDoListMongoDB(db *mongo.Client) *ToDoListMongoDB {
	return &ToDoListMongoDB{db}
}

func (r *ToDoListMongoDB) CreateTask(task models.Task) (int, error) {
	// Получение ссылки на коллекцию "tasks" в базе данных "testdb"
	collection := r.db.Database("testdb").Collection("tasks")

	// Создание счетчика, если его нет
	countersCollection := r.db.Database("testdb").Collection("counters")
	countersCollection.InsertOne(context.Background(), models.Counter{ID: "taskID", Sequence: 0})

	// Генерация нового _id на основе счетчика
	var newID int
	err := countersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": "taskID"}, bson.M{"$inc": bson.M{"seq": 1}}).Decode(&models.Counter{})
	if err != nil {
		log.Fatal(err)
	}

	// Получение текущего значения счетчика
	err = countersCollection.FindOne(context.Background(), bson.M{"_id": "taskID"}).Decode(&models.Counter{})
	if err != nil {
		log.Fatal(err)
	}

	// Вставка документа в MongoDB
	_, err = collection.InsertOne(context.Background(), task)
	if err != nil {
		log.Fatal(err)
	}
	return newID, nil
}
