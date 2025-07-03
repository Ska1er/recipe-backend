package repository

import (
	"database/sql"
	"fmt"

	"recipe/model"
)

type IngredientRepository interface {
	Get(page, pageSize int, search string) ([]model.Ingredient, int, error)
}

func NewIngredientRepository(db *sql.DB) IngredientRepository {
	return &IngredientPostgresRepository{db}
}

type IngredientPostgresRepository struct {
	db *sql.DB
}

func (repo *IngredientPostgresRepository) Get(page, pageSize int, search string) ([]model.Ingredient, int, error) {
	offset := (page - 1) * pageSize
	countQuery := `SELECT COUNT(*) FROM ingredients ing`
	query := `SELECT 
			ing.id,
			ing.name,
			ing.measure,
			ing.created_at,
			ing.updated_at
		FROM ingredients ing
		`
	var args []any = make([]any, 0, 3)
	if search != "" {
		query += ` WHERE ing.name ILIKE '%' || $1 || '%'`
		countQuery += ` WHERE ing.name ILIKE '%' || $1 || '%'`
		args = append(args, search)
	}

	query += fmt.Sprintf(` ORDER BY ing.created_at DESC
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

	ingredients := make([]model.Ingredient, 0, pageSize)
	for rows.Next() {
		var ingredient model.Ingredient
		rows.Scan(
			&ingredient.ID,
			&ingredient.Name,
			&ingredient.Measure,
			&ingredient.CreatedAt,
			&ingredient.UpdatedAt,
		)
		ingredients = append(ingredients, ingredient)
	}

	return ingredients, count, nil
}
