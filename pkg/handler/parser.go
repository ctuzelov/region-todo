package handler

import (
	"github.com/ctuzelov/region-todo/pkg/validator"
	"github.com/gin-gonic/gin"
)

func Parser(g *gin.Context, input *todoForm) error {
	if err := g.BindJSON(&input); err != nil {
		return err
	}

	if !validator.Valid(validator.NotBlank(input.Title), validator.MaxChars(input.Title, 200), validator.IsValidDate(input.ActiveAt)) {
		return errForm
	}
	return nil
}
