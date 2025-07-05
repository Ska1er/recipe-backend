package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"recipe/model"
)

type RecipeRepository interface {
	Save(recipe *model.Recipe) error
	GetByID(ID int) (*model.Recipe, error)
	Get(page, pageSize int, search string) ([]model.Recipe, int, error)
}

func NewRecipeRepository(db *sql.DB) RecipeRepository {
	return &RecipePostgresRepository{db}
}

type RecipePostgresRepository struct {
	db *sql.DB
}

func (repo *RecipePostgresRepository) GetByID(ID int) (*model.Recipe, error) {
	query := `SELECT 
			r.id,
			r.name,
			r.description,
			r.image,
			r.steps,
			r.is_custom,
			r.is_visible,
			r.cooking_time,
			r.difficulty,
			r.created_at,
			r.updated_at,
			ri.ingredient_id as ri_ingredient_id,
			ri.value as ri_value,
			i.id as i_id,
			i.name as i_iname,
			i.measure as i_measure
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON ri.recipe_id = r.id
		LEFT JOIN ingredients i ON i.id = ri.ingredient_id
		WHERE r.id = $1
		ORDER BY ri.sort ASC
		`
	rows, err := repo.db.Query(
		query,
		ID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	recipe := model.Recipe{
		Ingredients: make([]model.RecipeIngredient, 0, 20),
	}
	var steps []byte
	var hasValue bool = false
	for rows.Next() {
		hasValue = true
		var ri = model.RecipeIngredient{
			Ingredient: &model.Ingredient{},
		}

		err := rows.Scan(
			&recipe.ID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Image,
			&steps,
			&recipe.IsCustom,
			&recipe.IsVisible,
			&recipe.CookingTime,
			&recipe.Difficulty,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
			&ri.IngredientID,
			&ri.Value,
			&ri.Ingredient.ID,
			&ri.Ingredient.Name,
			&ri.Ingredient.Measure,
		)

		if err != nil {
			return nil, err
		}

		recipe.Ingredients = append(recipe.Ingredients, ri)

		if recipe.Steps == nil {
			var recipeSteps []model.RecipeStep
			err = json.Unmarshal(steps, &recipeSteps)
			if err != nil {
				return nil, err
			}

			recipe.Steps = recipeSteps
		}

	}

	if !hasValue {
		return nil, nil
	}

	return &recipe, nil
}

func (repo *RecipePostgresRepository) Save(recipe *model.Recipe) error {
	if recipe == nil {
		return errors.New("recipe is nil")
	}

	stepsJSON, err := json.Marshal(recipe.Steps)
	if err != nil {
		return err
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO recipes (name, description, image, steps, is_custom, is_visible, cooking_time, difficulty)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id, created_at, updated_at`
	err = tx.QueryRow(
		query,
		recipe.Name,
		recipe.Description,
		recipe.Image,
		stepsJSON,
		recipe.IsCustom,
		recipe.IsVisible,
		recipe.CookingTime,
		string(recipe.Difficulty),
	).Scan(&recipe.ID, &recipe.CreatedAt, &recipe.UpdatedAt)

	if err != nil {
		return err
	}

	for key, ingredient := range recipe.Ingredients {
		query := `INSERT INTO recipe_ingredients (sort, recipe_id, ingredient_id, value)
			VALUES($1, $2, $3, $4)`
		_, err = tx.Exec(
			query,
			key,
			recipe.ID,
			ingredient.IngredientID,
			ingredient.Value,
		)

		if err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

func (repo *RecipePostgresRepository) Get(page, pageSize int, search string) ([]model.Recipe, int, error) {
	offset := (page - 1) * pageSize
	countQuery := `SELECT COUNT(*) FROM recipes r WHERE r.is_visible = true`
	query := `SELECT 
			r.id,
			r.name,
			r.description,
			r.image,
			r.is_custom,
			r.is_visible,
			r.cooking_time,
			r.difficulty,
			r.created_at,
			r.updated_at
		FROM recipes r
		WHERE r.is_visible = true
		`
	var args []any = make([]any, 0, 3)
	if search != "" {
		query += ` AND r.name ILIKE '%' || $1 || '%'`
		countQuery += ` AND r.name ILIKE '%' || $1 || '%'`
		args = append(args, search)
	}

	query += fmt.Sprintf(` ORDER BY r.created_at DESC
		LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, pageSize, offset)

	rows, err := repo.db.Query(
		query,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var count int
	repo.db.QueryRow(countQuery, search).Scan(&count)

	recipes := make([]model.Recipe, 0, pageSize)
	for rows.Next() {
		var recipe model.Recipe
		rows.Scan(
			&recipe.ID,
			&recipe.Name,
			&recipe.Description,
			&recipe.Image,
			&recipe.IsCustom,
			&recipe.IsVisible,
			&recipe.CookingTime,
			&recipe.Difficulty,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)
		recipes = append(recipes, recipe)
	}

	return recipes, count, nil
}
