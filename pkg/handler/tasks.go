package handler

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/pkg/validator"
	"github.com/gin-gonic/gin"
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
	// * Получение данных из запроса
	var input todoForm

	// * Парсинг и валидация данных
	err := Parser(g, &input)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	// * Валидация полей задачи
	if !validator.Valid(validator.NotBlank(input.Title), validator.MaxChars(input.Title, 200), validator.IsValidDate(input.ActiveAt)) {
		ErrorResponse(g, errForm.Error())
		return
	}

	// * Преобразование строки времени в объект времени
	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	// * Создание и сохранение задачи
	// ? ActiveAt ковертируется с string в time.Timeчтобы записать в базу данных
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

	// * Отправка успешного статуса
	g.Status(http.StatusNoContent)
}

// getTasksByStatus gets tasks by status.
//
//	@Summary		Get tasks by status
//	@Description	Get tasks by the specified status or default to 'active'
//	@Tags			Tasks
//	@Param			status	query	string	false	"Task status (default: active)"
//	@Success		200		{array}	todoForm
//	@Router			/api/todo-list/tasks [get]
func (h *Handler) getTasksByStatus(g *gin.Context) {
	// * Получение параметра статуса из запроса, по умолчанию "active".
	status := g.DefaultQuery("status", "active")

	// * Чтение задач из сервиса и обработка возможной ошибки.
	tasks, err := h.service.ReadTasks(status)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	currentDate := time.Now()

	// * Сортировка задач по дате создания (по возрастанию).
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	var response []todoForm

	// * Обработка каждой задачи для формирования ответа.
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

	// ? Возвращение ответа в JSON формате, с учетом возможной пустоты.
	if len(response) == 0 {
		g.JSON(http.StatusOK, []todoForm{})
		return
	}
	g.JSON(http.StatusOK, response)
}

// deleteTaskByID deletes a task by ID.
//
//	@Summary		Delete a task by ID
//	@Description	Delete a task by the specified ID
//	@Tags			Tasks
//	@Produce		json
//	@Param			id	path	int	true	"Task ID"
//	@Success		204	"No Content"
//	@Failure		404	"Not Found"	Task	not	found
//	@Router			/api/todo-list/tasks/{id} [delete]
func (h *Handler) deleteTaskByID(g *gin.Context) {
	// * Получение идентификатора задачи из параметров URL
	idStr := g.Param("id")

	// * Преобразование строки идентификатора в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "неверный тип")
		return
	}

	// * Вызов сервисной функции для удаления задачи по идентификатору
	err = h.service.DeleteTask(id)

	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}

// updateTaskStatusByID updates a task status by ID.
//
//	@Summary		Update a task status by ID
//	@Description	Updates the status of a task by ID.
//	@Tags			Tasks
//	@Produce		json
//	@Param			id	path	int	true	"The ID of the task to update"
//	@Success		204	"No Content"
//	@Failure		404	"Not Found"		Task	not	found
//	@Router			/api/todo-list/tasks/{id}/done [put]
func (h *Handler) updateTaskStatusByID(g *gin.Context) {
	// * Получение ID задачи из параметра URL
	idStr := g.Param("id")

	// * Преобразование строки ID в целое число
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "invalid type")
		return
	}

	// * Обновление статуса задачи через сервис
	err = h.service.UpdateTaskStatus(id)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}
	g.Status(http.StatusNoContent)
}

// updateTaskByID updates a task by ID.
//
//	@Summary		Update a task by ID
//	@Description	Update a task with the provided input
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int			true	"Task ID"
//	@Param			input	body	todoForm	true	"Task input"
//	@Success		204		"No Content"
//	@Failure		404		"Not Found"
//	@Router			/api/todo-list/tasks/{id} [put]
func (h *Handler) updateTaskByID(g *gin.Context) {
	var input todoForm

	// * Получение ID из параметра маршрута.
	idStr := g.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		ErrorResponse(g, "invalid type")
		return
	}

	// * Парсинг и валидация входных данных.
	err = Parser(g, &input)
	if err != nil || !validator.Valid(validator.NotBlank(input.Title), validator.MaxChars(input.Title, 200), validator.IsValidDate(input.ActiveAt)) {
		ErrorResponse(g, errForm.Error())
		return
	}

	// * Преобразование строки даты в формат time.Time.
	convertedtime, _ := time.Parse("2006-01-02", input.ActiveAt)

	// * Создание и сохранение задачи
	// ? ActiveAt ковертируется с string в time.Timeчтобы записать в базу данных
	validInput := models.Task{
		Title:    input.Title,
		ActiveAt: convertedtime,
	}

	// * Обновление задачи через сервис и обработка ошибок.
	err = h.service.ToDoTasks.UpdateTask(id, validInput)
	if err != nil {
		ErrorResponse(g, err.Error())
		return
	}

	g.Status(http.StatusNoContent)
}
