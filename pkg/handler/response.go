package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var errForm = errors.New("follow next rules: title length <= 200 and not blank, ativeAt format 2023-01-01")
var errDuplicate = errors.New("a task with such values already exists")

func ErrorResponse(g *gin.Context, message string) {
	logrus.Error(message)
	g.Status(http.StatusNotFound)
}
