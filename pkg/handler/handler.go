package handler

import (
	"github.com/ctuzelov/region-todo/pkg/service"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ctuzelov/region-todo/docs"
)

// * Объявление структуры Handler, которая будет обрабатывать запросы
type Handler struct {
	service *service.Service
}

// * Создание нового экземпляра Handler с передачей сервиса в качестве зависимости
func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

// * Инициализация маршрутов Gin и добавление обработчиков для различных запросов
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// * Добавление маршрута для обслуживания Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// * Группирование маршрутов для задач
	tasks := router.Group("/api/todo-list/tasks")

	// * Привязка обработчиков к соответствующим HTTP-методам и маршрутам
	tasks.POST("/", h.createTask)
	tasks.GET("/", h.getTasksByStatus)
	tasks.PUT("/:id/done", h.updateTaskStatusByID)
	tasks.PUT("/:id", h.updateTaskByID)
	tasks.DELETE("/:id", h.deleteTaskByID)

	return router
}
