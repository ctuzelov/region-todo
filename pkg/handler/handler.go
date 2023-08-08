package handler

import (
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	tasks := router.Group("/api/todo-list/tasks")

	tasks.POST("/", h.createTask)
	tasks.GET("/", h.getTasksByStatus)
	tasks.PUT("/:id/done", h.updateTaskStatusByID)
	tasks.PUT("/:id", h.updateTaskByID)
	tasks.DELETE("/:id", h.deleteTaskByID)

	return router
}
