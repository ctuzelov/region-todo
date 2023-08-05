package handler

import (
	"net/http"

	"github.com/ctuzelov/region-todo/pkg/models"
	"github.com/ctuzelov/region-todo/util"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(g *gin.Context) {
	var input models.Task
	if err := g.BindJSON(&input); err != nil {
		NewErrorResponse(g, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.ToDoList.CreateTask(input)
	if err != nil {
		NewErrorResponse(g, http.StatusUnauthorized, err.Error())
		return
	}

	g.JSON(http.StatusOK, models.Task{
		ID:       util.IntToByteArray(id),
		Title:    input.Title,
		ActiveAt: input.ActiveAt,
	})
}

func (h *Handler) getListByID(g *gin.Context) {

}

func (h *Handler) updateListByID(g *gin.Context) {

}

func (h *Handler) deleteListByID(g *gin.Context) {

}
