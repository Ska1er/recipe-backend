package internal

import (
	"recipe/model"
)

type CreateRecipeRequest struct {
	Name        string                           `json:"name" binding:"required,max=100"`
	Description *string                          `json:"description" binding:"required"`
	Image       string                           `json:"image" binding:"required,max=200"`
	CookingTime int                              `json:"cookingTime" binding:"required"`
	Difficulty  model.RecipeDifficulty           `json:"difficulty" binding:"required"`
	Steps       []*CreateRecipeStepRequest       `json:"steps" binding:"required"`
	Ingredients []*CreateRecipeIngredientRequest `json:"ingredients"`
}

func (req *CreateRecipeRequest) ToEntity() *model.Recipe {
	steps := make([]model.RecipeStep, 0, len(req.Steps))
	for _, step := range req.Steps {
		steps = append(steps, *step.ToEntity())
	}

	ingredients := make([]model.RecipeIngredient, 0, len(req.Ingredients))
	for _, ingredient := range req.Ingredients {
		ingredients = append(ingredients, *ingredient.ToEntity())
	}

	return &model.Recipe{
		Name:        req.Name,
		Description: ToNullString(req.Description),
		Image:       req.Image,
		Steps:       steps,
		Ingredients: ingredients,
		CookingTime: req.CookingTime,
		Difficulty:  req.Difficulty,
	}
}

type CreateRecipeStepRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Description string `json:"description" binding:"required,max=500"`
	Image       string `json:"image" binding:"requried,max=200"`
}

func (req *CreateRecipeStepRequest) ToEntity() *model.RecipeStep {
	return &model.RecipeStep{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
	}
}

type CreateRecipeIngredientRequest struct {
	IngredientID int     `json:"id" binding:"required"`
	Value        float64 `json:"value" binding:"required,min=0"`
}

func (req *CreateRecipeIngredientRequest) ToEntity() *model.RecipeIngredient {
	return &model.RecipeIngredient{
		IngredientID: req.IngredientID,
		Value:        req.Value,
	}
}
