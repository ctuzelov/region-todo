package handler

import (
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	tasks := router.Group("/api/todo-list/tasks")

	tasks.POST("/", h.createTask)
	tasks.GET("/", h.getTasksByStatus)
	tasks.PUT("/:id/done", h.updateTaskStatusByID)
	tasks.PUT("/:id", h.updateTaskByID)
	tasks.DELETE("/:id", h.deleteTaskByID)

	return router
}
