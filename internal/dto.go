package internal

import (
	"recipe/model"
)

type RecipeDTO struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description *string                `json:"description"`
	Image       string                 `json:"image"`
	CookingTime int                    `json:"cooking_time"`
	Difficulty  model.RecipeDifficulty `json:"difficulty"`
	Steps       []model.RecipeStep     `json:"steps"`
	Ingredients []IngredientDTO        `json:"ingredients"`
}

func NewRecipeFromEntity(recipe model.Recipe) *RecipeDTO {
	var description *string = nil
	if recipe.Description.Valid {
		description = &recipe.Description.String
	}
	ingredients := make([]IngredientDTO, 0, len(recipe.Ingredients))
	for _, ing := range recipe.Ingredients {
		ingredients = append(ingredients, IngredientDTO{
			ID:      ing.Ingredient.ID,
			Name:    ing.Ingredient.Name,
			Measure: ing.Ingredient.Measure,
			Value:   ing.Value,
		})
	}

	return &RecipeDTO{
		ID:          recipe.ID,
		Name:        recipe.Name,
		Description: description,
		Image:       recipe.Image,
		CookingTime: recipe.CookingTime,
		Difficulty:  recipe.Difficulty,
		Steps:       recipe.Steps,
		Ingredients: ingredients,
	}
}

type IngredientDTO struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Measure string  `json:"measure"`
	Value   float64 `json:"value"`
}
