package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"recipe/repository"

	"github.com/gin-gonic/gin"
)

type IngredientHandler struct {
	ingredientRepository repository.IngredientRepository
}

func NewIngredientHandler(ingredientRepository repository.IngredientRepository) *IngredientHandler {
	return &IngredientHandler{
		ingredientRepository: ingredientRepository,
	}
}

func (handler *IngredientHandler) ListHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	search := c.DefaultQuery("search", "")

	ingredients, count, err := handler.ingredientRepository.Get(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed get ingredients: %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ingredients": ingredients,
		"totalPage":   math.Ceil(float64(count / pageSize)),
		"total":       count,
	})
}
