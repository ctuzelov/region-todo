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

// createTask godoc
// @Summary Create todo task
// @Tags tasks
// @Description Create a new task
// @ID create-task
// @Accept json
// @Produce json
// @Param input body todoForm true "Task details to be created"
// @Success 204 {}
// @Router /api/todo-list/tasks [post]

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

//	getTasksByStatus godoc
//	@Summary		Get tasks by status
//	@Tags			tasks
//	@Description	Retrieve tasks based on their status
//	@ID				get-tasks-by-status
//	@Accept			json
//	@Produce		json
//	@Param			status	query		string			false	"Task status ('active' or 'done') (default: 'active')"
//	@Success		200		{array}		string		"List of tasks with their titles and activeAt dates"
//	@Failure		404		{object}	errForm	"Not Found"
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

	if len(response) == 0 {
		g.JSON(http.StatusOK, []todoForm{})
		return
	}

	g.JSON(http.StatusOK, response)
}

//	deleteTaskByID  godoc
//	@Summary		Delete task by ID
//	@Tags			tasks
//	@Description	Delete a task by its ID
//	@ID				delete-task-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	integer	true	"Task ID to delete"
//	@Success		204	"Task deleted successfully"
//	@Failure		404	{object}	errForm	"Not Found"
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

//	updateTaskStatusByID godoc
//	@Summary		Update task status by ID
//	@Tags			tasks
//	@Description	Update the status of a task by its ID
//	@ID				update-task-status-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id	path	integer	true	"Task ID to update status"
//	@Success		204	"Task status updated successfully"
//	@Failure		404	{object}	errForm	"Not Found"
//	@Router			/api/todo-list/tasks/{id}/done [put]

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

// updateTaskByID	godoc
//	@Summary		Update task by ID
//	@Tags			tasks
//	@Description	Update a task by its ID
//	@ID				update-task-by-id
//	@Accept			json
//	@Produce		json
//	@Param			id		path	integer		true	"Task ID to update"
//	@Param			input	body	todoForm	true	"Task details to be updated"
//	@Success		204		"Task updated successfully"
//	@Failure		404		{object}	errForm	"Not Found"
//	@Router			/api/todo-list/tasks/{id} [put]

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
