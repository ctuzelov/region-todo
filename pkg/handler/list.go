package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type todoForm struct {
	Title    string `bson:"title" json:"title"`
	ActiveAt string `bson:"activeAt" json:"activeAt"`
}

func (h *Handler) createList(g *gin.Context) {
	var input todoForm

	err := Parser(g, &input)
	if err != nil {
		NewErrorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	validInput := models.Task{
		Title:    input.Title,
		ActiveAt: convertedtime,
	}

	_, err = h.service.ToDoList.CreateTask(validInput)
	if errors.Is(err, errDuplicate) {
		NewErrorResponse(g, http.StatusUnauthorized, errDuplicate.Error())
		return
	} else if err != nil {
		NewErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}

func (h *Handler) getListByID(g *gin.Context) {
	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}

	todo, err := h.service.ToDoList.ReadTask(id)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			NewErrorResponse(g, http.StatusNotFound, err.Error())
			return
		}
		g.JSON(http.StatusNotFound, gin.H{"error": "Database error"})
		return
	}

	g.JSON(http.StatusOK, todoForm{
		Title:    todo.Title,
		ActiveAt: todo.ActiveAt.Format("2006-01-02"),
	})
}

func (h *Handler) deleteListByID(g *gin.Context) {
	// Get the task ID from the URL parameters
	idStr := g.Param("id")

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.service.Delete(id)

	if err != nil {
		NewErrorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}

func (h *Handler) updatTaskStatusByID(g *gin.Context) {
	// Get the task ID from the URL parameters
	idStr := g.Param("id")

	// Convert the ID string to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	err = h.service.UpdateStatus(id)
	if err != nil {
		NewErrorResponse(g, http.StatusNotFound, err.Error())
		return
	}
	g.Status(http.StatusNoContent)
}

func (h *Handler) updatTaskByID(g *gin.Context) {
	var input todoForm
	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		NewErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}

	err = Parser(g, &input)
	if err != nil {
		NewErrorResponse(g, http.StatusNotFound, err.Error())
		return
	}

	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	validInput := models.Task{
		Title:    input.Title,
		ActiveAt: convertedtime,
	}

	err = h.service.ToDoList.UpdateTask(id, validInput)
	if err != nil {
		NewErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}
