package handler

import (
	"net/http"

	"github.com/ctuzelov/region-todo/pkg/validator"
	"github.com/gin-gonic/gin"
)

func Parser(g *gin.Context, input *todoForm) error {
	if err := g.BindJSON(&input); err != nil {
		NewErrorResponse(g, http.StatusBadRequest, err.Error())
		return err
	}

	if !validator.Valid(validator.NotBlank(input.Title), validator.MaxChars(input.Title, 201), validator.IsValidDate(input.ActiveAt)) {
		NewErrorResponse(g, http.StatusUnauthorized, errForm.Error())
		return errForm
	}
	return nil
}
