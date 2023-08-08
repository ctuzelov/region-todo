package handler

import (
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ctuzelov/region-todo/docs"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

// tasksRouter defines operations for managing tasks.
//	@Summary		Manage tasks
//	@Description	Operations for managing tasks
//	@Tags			Tasks
//	@Router			/api/todo-list/tasks/{id} [delete]
//	@Router			/api/todo-list/tasks [post]
//	@Router			/api/todo-list/tasks [get]
//	@Router			/api/todo-list/tasks/{id} [put]
//	@Router			/api/todo-list/tasks/{id}/done [put]

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	tasks := router.Group("/api/todo-list/tasks")

	tasks.POST("/", h.createTask)
	tasks.GET("/", h.getTasksByStatus)
	tasks.PUT("/:id/done", h.updateTaskStatusByID)
	tasks.PUT("/:id", h.updateTaskByID)
	tasks.DELETE("/:id", h.deleteTaskByID)

	return router
}
