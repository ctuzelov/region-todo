package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

var errForm = errors.New("follow next rules: title length <= 200, ativeAt format 2023-01-01")
var errDuplicate = errors.New("a task with such values already exists")

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
