package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/gin-swagger"
)

type todoForm struct {
	Title    string `bson:"title" json:"title"`
	ActiveAt string `bson:"activeAt" json:"activeAt"`
}

// createTask creates a new task.
//
//	@Summary		Create a new task
//	@Description	Create a new task with the provided input
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			input	body	todoForm	true	"Task input"
//	@Success		204		"No Content"
//	@Failure		404		"Not Found"
//	@Router			/api/todo-list/tasks [post]
func (h *Handler) createTask(g *gin.Context) {
	var input todoForm

	err := Parser(g, &input)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	validInput := models.Task{
		Title:    input.Title,
		ActiveAt: convertedtime,
	}

	_, err = h.service.ToDoTasks.CreateTask(validInput)
	if errors.Is(err, errDuplicate) {
		ErrorResponse(g, errDuplicate.Error())
		return
	} else if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}

// getTasksByStatus gets tasks by status.
//	@Summary		Get tasks by status
//	@Description	Get tasks by the specified status or default to 'active'
//	@Tags			Tasks
//	@Produce		json
//	@Param			status	query		string	false	"Task status (default: active)"
//	@Success		200		{object}	todoForm
//	@Failure		400,	404			{object}	ErrorResponse
//	@Router			/api/todo-list/tasks [get]

func (h *Handler) getTasksByStatus(g *gin.Context) {
	status := g.DefaultQuery("status", "active")

	tasks, err := h.service.ReadTasks(status)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	currentDate := time.Now()

	// Sort the tasks based on CreatedAt date (ascending order)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	var response []todoForm

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ActiveAt.Weekday() == time.Saturday || tasks[i].ActiveAt.Weekday() == time.Sunday {
			tasks[i].Title = fmt.Sprintf("ВЫХОДНОЙ - %s", tasks[i].Title)
		}

		if tasks[i].Status == "active" && tasks[i].ActiveAt.Before(currentDate) {
			response = append(response, todoForm{
				Title:    tasks[i].Title,
				ActiveAt: tasks[i].ActiveAt.Format("2006-01-02"),
			})
		} else if tasks[i].Status == "done" {
			response = append(response, todoForm{
				Title:    tasks[i].Title,
				ActiveAt: tasks[i].ActiveAt.Format("2006-01-02"),
			})
		}
	}

	g.JSON(http.StatusOK, response)
}

//	@Summary		Delete a task by ID
//	@Description	Delete a task by the specified ID
//	@Tags			Tasks
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204	"No Content"
//	@Failure		404	"Not Found"
//	@Router			/api/todo-list/tasks/{id} [delete]

func (h *Handler) deleteTaskByID(g *gin.Context) {
	// Get the task ID from the URL parameters
	idStr := g.Param("id")

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "invalid type")
		return
	}

	err = h.service.DeleteTask(id)

	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}

func (h *Handler) updateTaskStatusByID(g *gin.Context) {
	// Get the task ID from the URL parameters
	idStr := g.Param("id")

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "invalid type")
		return
	}

	err = h.service.UpdateTaskStatus(id)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}
	g.Status(http.StatusNoContent)
}

// updateTaskByID updates a task by ID.
//	@Summary		Update a task by ID
//	@Description	Update a task with the provided input
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			input	body	todoForm	true	"Task input"
//	@Success		204		"No Content"
//	@Failure		404		"Not Found"
//	@Router			/api/todo-list/tasks/:id [put]

func (h *Handler) updateTaskByID(g *gin.Context) {
	var input todoForm
	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "invalid type")
		return
	}

	err = Parser(g, &input)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	validInput := models.Task{
		Title:    input.Title,
		ActiveAt: convertedtime,
	}

	err = h.service.ToDoTasks.UpdateTask(id, validInput)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}
