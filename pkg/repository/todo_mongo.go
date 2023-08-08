package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ctuzelov/region-todo/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ToDoTasksMongoDB struct {
	db *mongo.Client
}

const newStatus = "done"

func NewToDoTasksMongoDB(db *mongo.Client) *ToDoTasksMongoDB {
	return &ToDoTasksMongoDB{db}
}

func (r *ToDoTasksMongoDB) CreateTask(task models.Task) (int, error) {
	// * Получение ссылки на коллекцию "tasks" в базе данных "taskdb"
	tasksCollection := r.db.Database("taskdb").Collection("tasks")

	// * Создание счетчика, если его нет
	countersCollection := r.db.Database("taskdb").Collection("counters")
	countersCollection.InsertOne(context.Background(), models.Counter{ID: "taskID", Sequence: 0})

	// * Генерация нового _id на основе счетчика
	var counter models.Counter
	err := countersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": "taskID"}, bson.M{"$inc": bson.M{"sequence": 1}}).Decode(&counter)
	if err != nil {
		return 0, err
	}

	// * Получение текущего значения счетчика
	err = countersCollection.FindOne(context.Background(), bson.M{"_id": "taskID"}).Decode(&models.Counter{})
	if err != nil {
		return 0, err
	}

	todo := models.Task{
		ID:        counter.Sequence + 1,
		Title:     task.Title,
		Status:    "active",
		ActiveAt:  task.ActiveAt,
		CreatedAt: time.Now(),
	}

	// * Проверяем наличие дубликата
	filter := bson.M{"title": task.Title, "activeAt": task.ActiveAt}
	err = tasksCollection.FindOne(context.Background(), filter).Decode(&todo)
	if err == nil {
		countersCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": "taskID"}, bson.M{"$inc": bson.M{"sequence": -1}})
		return 0, errors.New("a task with such values already exists")
	} else if err != mongo.ErrNoDocuments {
		return 0, err
	}

	// ! Вставка документа в MongoDB
	_, err = tasksCollection.InsertOne(context.Background(), todo)
	if err != nil {
		return 0, err
	}
	return counter.Sequence + 1, nil
}

func (r *ToDoTasksMongoDB) ReadTasks(status string) ([]models.Task, error) {
	// * Подключение к коллекции "tasks" в базе данных "taskdb"
	collection := r.db.Database("taskdb").Collection("tasks")

	filter := bson.M{"status": status}

	// * Поиск документа по указанному id
	var todo []models.Task
	all, err := collection.Find(context.Background(), filter)

	// * Итерация по найденным задачам
	for all.Next(context.Background()) {
		var task models.Task
		if err := all.Decode(&task); err != nil {
			return []models.Task{}, err
		}
		todo = append(todo, task)
	}

	if err := all.Err(); err != nil {
		return []models.Task{}, err
	}
	defer all.Close(context.Background())

	return todo, err
}

func (r *ToDoTasksMongoDB) DeleteTask(id int) error {
	// * Подключение к коллекции "tasks" в базе данных "taskdb"
	tasksCollection := r.db.Database("taskdb").Collection("tasks")

	// * Создание фильтра для поиска задачи по ID
	filter := bson.M{"_id": id}

	// ! Удаление задачи из коллекции
	res, err := tasksCollection.DeleteOne(context.Background(), filter)
	if res.DeletedCount == 0 {
		return fmt.Errorf("no object with the given id = %d", id)
	}
	return err
}

func (r *ToDoTasksMongoDB) UpdateTaskStatus(id int) error {
	// * Подключение к коллекции "tasks" в базе данных "taskdb"
	tasksCollection := r.db.Database("taskdb").Collection("tasks")

	// * Создание фильтра для поиска задачи по ID
	filter := bson.M{"_id": id}

	// ! Обновление поле "статус" в задаче
	update := bson.M{"$set": bson.M{"status": newStatus}}
	res, err := tasksCollection.UpdateOne(context.Background(), filter, update)
	if res.MatchedCount == 0 {
		return fmt.Errorf("no object with the given id = %d", id)
	}
	return err
}

func (r *ToDoTasksMongoDB) UpdateTask(id int, task models.Task) error {
	// * Подключение к коллекции "tasks" в базе данных "taskdb"
	tasksCollection := r.db.Database("taskdb").Collection("tasks")

	// * Создание фильтра для поиска задачи по ID
	filter := bson.M{"_id": id}

	// * Создаем обновление только для полей Title и ActiveAt
	update := bson.M{"$set": bson.M{"title": task.Title, "activeAt": task.ActiveAt}}

	// ! Обновление полей Title и ActiveAt в задаче
	res, err := tasksCollection.UpdateOne(context.Background(), filter, update)
	if res.MatchedCount == 0 {
		return fmt.Errorf("no object with the given id = %d", id)
	}
	return err
}
