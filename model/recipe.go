package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type RecipeDifficulty string

const (
	Easy   RecipeDifficulty = "easy"
	Medium RecipeDifficulty = "medium"
	Hard   RecipeDifficulty = "hard"
)

// Implement sql.Scanner interface
func (d *RecipeDifficulty) Scan(value any) error {
	if value == nil {
		*d = "" // or set a default value
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*d = RecipeDifficulty(string(v))
	case string:
		*d = RecipeDifficulty(v)
	default:
		return fmt.Errorf("unsupported type for Difficulty: %T", value)
	}

	// Optional: validate the value
	switch *d {
	case Easy, Medium, Hard:
		return nil
	default:
		return fmt.Errorf("invalid difficulty value: %s", *d)
	}
}

func (d *RecipeDifficulty) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*d = RecipeDifficulty(strings.ToLower(s))
	if *d != Easy && *d != Medium && *d != Hard {
		return fmt.Errorf("invalid difficulty: %s", s)
	}
	return nil
}

type Recipe struct {
	ID          int                `json:"id"`
	Name        string             `json:"name"`
	Description sql.NullString     `json:"description"`
	Image       string             `json:"image"`
	IsCustom    bool               `json:"is_custom"`
	IsVisible   bool               `json:"is_visible"`
	CookingTime int                `json:"cooking_time"`
	Difficulty  RecipeDifficulty   `json:"difficulty"`
	Steps       []RecipeStep       `json:"steps"`
	Ingredients []RecipeIngredient `json:"ingredients"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type RecipeStep struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type RecipeIngredient struct {
	IngredientID int         `json:"ingredient_id"`
	Ingredient   *Ingredient `json:"ingredient"`
	Value        float64     `json:"value"`
}

type Ingredient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Measure   string    `json:"measure"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
