package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

	"recipe/internal"
	"recipe/repository"

	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	recipeRepository repository.RecipeRepository
}

func NewRecipeHandler(recipeRepository repository.RecipeRepository) *RecipeHandler {
	return &RecipeHandler{recipeRepository}
}

func (handler *RecipeHandler) CreateHandler(c *gin.Context) {
	var request internal.CreateRecipeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	recipe := request.ToEntity()
	recipe.IsCustom = true
	recipe.IsVisible = false

	err := handler.recipeRepository.Save(recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed save recipe: %s", err.Error()),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (handler *RecipeHandler) GetHandler(c *gin.Context) {
	ID, _ := strconv.Atoi(c.Param("id"))
	recipe, err := handler.recipeRepository.GetByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed get recipe: %s", err.Error()),
		})
		return
	}

	if recipe == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Recipe not found",
		})
		return
	}

	dto := internal.NewRecipeFromEntity(*recipe)
	c.JSON(http.StatusOK, dto)
}

func (handler *RecipeHandler) ListHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "1"))
	search := c.DefaultQuery("search", "")

	recipes, count, err := handler.recipeRepository.Get(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed get recipes: %s", err.Error()),
		})
		return
	}

	dto := make([]internal.RecipeDTO, 0, len(recipes))
	for _, recipe := range recipes {
		dto = append(dto, *internal.NewRecipeFromEntity(recipe))
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes":   dto,
		"totalPage": math.Ceil(float64(count / pageSize)),
		"total":     count,
	})
}
