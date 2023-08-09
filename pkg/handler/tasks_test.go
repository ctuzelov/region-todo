package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/service"
	mock_service "github.com/ctuzelov/region-todo/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/magiconair/properties/assert"
)

func TestHandler_createTask(t *testing.T) {
	convertedtime, _ := time.Parse("2006-01-02", "2023-08-05")

	testTable := []struct {
		name               string
		inputBody          string
		inputTodo          models.Task
		mockBehavior       func(s *mock_service.MockToDoTasks, todo models.Task)
		expectedStatusCode int
	}{
		{
			name:      "NO CONTENT",
			inputBody: `{"title":"купить машину","activeAt":"2023-08-05"}`,
			inputTodo: models.Task{
				Title:    "купить машину",
				ActiveAt: convertedtime,
			},
			mockBehavior: func(s *mock_service.MockToDoTasks, todo models.Task) {
				s.EXPECT().CreateTask(todo).Return(1, nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:      "Error",
			inputBody: `{"title":"","activeAt":"2023-08-05"}`,
			inputTodo: models.Task{
				Title:    "",
				ActiveAt: convertedtime,
			},
			mockBehavior:       func(s *mock_service.MockToDoTasks, todo models.Task) {},
			expectedStatusCode: 404,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			taskcreate := mock_service.NewMockToDoTasks(c)
			testCase.mockBehavior(taskcreate, testCase.inputTodo)

			service := &service.Service{ToDoTasks: taskcreate}
			handler := NewHandler(service)

			r := gin.New()

			r.POST("/", handler.createTask)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_updateTaskStatusByID(t *testing.T) {
	testTable := []struct {
		name               string
		taskID             int
		newStatus          string
		mockBehavior       func(s *mock_service.MockToDoTasks, taskID int)
		expectedStatusCode int
	}{
		{
			name:      "Success",
			taskID:    1,
			newStatus: "completed",
			mockBehavior: func(s *mock_service.MockToDoTasks, taskID int) {
				s.EXPECT().UpdateTaskStatus(taskID).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:               "Task not found",
			taskID:             -1,
			newStatus:          "completed",
			mockBehavior:       func(s *mock_service.MockToDoTasks, taskID int) {},
			expectedStatusCode: 404,
		},
		// Add more test cases for different scenarios
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mock_service.NewMockToDoTasks(c)
			testCase.mockBehavior(taskService, testCase.taskID)

			service := &service.Service{ToDoTasks: taskService}
			handler := NewHandler(service)

			r := gin.New()
			r.PUT("/:id", handler.updateTaskStatusByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/"+strconv.Itoa(testCase.taskID), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_updateTaskByID(t *testing.T) {
	convertedTime, _ := time.Parse("2006-01-02", "2023-08-05")

	testTable := []struct {
		name               string
		taskID             int
		inputBody          string
		inputTask          models.Task
		mockBehavior       func(s *mock_service.MockToDoTasks, taskID int, task models.Task)
		expectedStatusCode int
	}{
		{
			name:      "Success",
			taskID:    0,
			inputBody: `{"title":"купить книгу","activeAt":"2023-08-05"}`,
			inputTask: models.Task{
				ID:       0,
				Title:    "купить книгу",
				ActiveAt: convertedTime,
			},
			mockBehavior: func(s *mock_service.MockToDoTasks, taskID int, task models.Task) {
				s.EXPECT().UpdateTask(taskID, task).Return(nil)
			},
			expectedStatusCode: 204,
		},
		{
			name:      "Task not found",
			taskID:    -1,
			inputBody: `{"title":"купить книгу","activeAt":"2023-08-05"}`,
			inputTask: models.Task{
				ID:       -1,
				Title:    "купить книгу",
				ActiveAt: convertedTime,
			},
			mockBehavior:       func(s *mock_service.MockToDoTasks, taskID int, task models.Task) {},
			expectedStatusCode: 404,
		},
		// Add more test cases for different scenarios
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mock_service.NewMockToDoTasks(c)
			testCase.mockBehavior(taskService, testCase.taskID, testCase.inputTask)

			service := &service.Service{ToDoTasks: taskService}
			handler := NewHandler(service)

			r := gin.New()
			r.PUT("/:id", handler.updateTaskByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/"+strconv.Itoa(testCase.taskID), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_getTasksByStatus(t *testing.T) {
	testTable := []struct {
		name                 string
		queryParams          map[string]string
		mockBehavior         func(s *mock_service.MockToDoTasks, status string) ([]models.Task, error)
		expectedResponseBody string
		expectedStatusCode   int
	}{
		{
			name: "Success",
			queryParams: map[string]string{
				"status": "active",
			},
			mockBehavior: func(s *mock_service.MockToDoTasks, status string) ([]models.Task, error) {
				tasks := []models.Task{
					{
						Title:     "Task 1",
						Status:    "active",
						ActiveAt:  time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC),
						CreatedAt: time.Date(2023, 8, 9, 0, 0, 0, 0, time.UTC),
					},
					{
						Title:     "Task 2",
						Status:    "active",
						ActiveAt:  time.Date(2023, 8, 15, 0, 0, 0, 0, time.UTC),
						CreatedAt: time.Date(2023, 8, 10, 0, 0, 0, 0, time.UTC),
					},
				}
				return tasks, nil
			},
			expectedResponseBody: `[{"title":"ВЫХОДНОЙ - Task 1","activeAt":"2023-08-05"},{"title":"Task 2","activeAt":"2023-08-15"}]`,
			expectedStatusCode:   200,
		},
		{
			name: "Error",
			queryParams: map[string]string{
				"status": "done",
			},
			mockBehavior: func(s *mock_service.MockToDoTasks, status string) ([]models.Task, error) {
				return []models.Task{}, nil
			},
			expectedResponseBody: `[]`,
			expectedStatusCode:   200,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mock_service.NewMockToDoTasks(c)
			taskService.EXPECT().ReadTasks(testCase.queryParams["status"]).Return(testCase.mockBehavior(taskService, testCase.queryParams["status"]))

			service := &service.Service{ToDoTasks: taskService}
			handler := NewHandler(service)

			r := gin.New()
			r.GET("/", handler.getTasksByStatus)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/?status=%s", testCase.queryParams["status"]), nil)

			query := req.URL.Query()
			for key, value := range testCase.queryParams {
				query.Add(key, value)
			}
			req.URL.RawQuery = query.Encode()

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_deleteTaskByID(t *testing.T) {
	testTable := []struct {
		name               string
		taskID             int
		mockBehavior       func(s *mock_service.MockToDoTasks, taskID int) error
		expectedStatusCode int
	}{
		{
			name:   "Success",
			taskID: 1,
			mockBehavior: func(s *mock_service.MockToDoTasks, taskID int) error {
				return nil
			},
			expectedStatusCode: 204,
		},
		{
			name:   "Task not found",
			taskID: 0,
			mockBehavior: func(s *mock_service.MockToDoTasks, taskID int) error {
				return errors.New("invalid type")
			},
			expectedStatusCode: 404,
		},
		// Add more test cases for different scenarios
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			taskService := mock_service.NewMockToDoTasks(c)

			// Set up the expectation for the UpdateTaskStatus method
			taskService.EXPECT().DeleteTask(testCase.taskID).Return(testCase.mockBehavior(taskService, testCase.taskID))

			service := &service.Service{ToDoTasks: taskService}
			handler := NewHandler(service)

			r := gin.New()
			r.DELETE("/:id", handler.deleteTaskByID)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/"+strconv.Itoa(testCase.taskID), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
