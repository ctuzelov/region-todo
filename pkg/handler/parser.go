package handler

import (
	"github.com/gin-gonic/gin"
)

func Parser(g *gin.Context, input *todoForm) error {
	if err := g.BindJSON(&input); err != nil {
		return err
	}

	return nil
}
