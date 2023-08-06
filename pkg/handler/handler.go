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

	tasks.POST("/", h.createList)
	tasks.GET("/:id", h.getListByID)
	tasks.PUT("/:id/done", h.updatTaskStatusByID)
	tasks.PUT("/:id/", h.updatTaskByID)
	tasks.DELETE("/:id", h.deleteListByID)

	return router
}
