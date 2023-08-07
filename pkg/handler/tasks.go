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
)

type todoForm struct {
	Title    string `bson:"title" json:"title"`
	ActiveAt string `bson:"activeAt" json:"activeAt"`
}

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

func (h *Handler) getTasksByStatus(g *gin.Context) {
	status := g.DefaultQuery("status", "active")

	tasks, err := h.service.ReadTasks(status)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	// Sort the tasks based on CreatedAt date (ascending order)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	response := make([]todoForm, len(tasks))

	for i := 0; i < len(tasks); i++ {
		if tasks[i].ActiveAt.Weekday() == time.Saturday || tasks[i].ActiveAt.Weekday() == time.Sunday {
			tasks[i].Title = fmt.Sprintf("ВЫХОДНОЙ - %s", tasks[i].Title)
		}
		response[i] = todoForm{
			Title:    tasks[i].Title,
			ActiveAt: tasks[i].ActiveAt.Format("2006-01-02"),
		}
	}

	g.JSON(http.StatusOK, response)
}

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
