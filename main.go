package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"

	"recipe/handler"
	"recipe/internal"
	"recipe/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		panic("Env variable DB_URL is not set")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic("Cannot connect to db")
	}

	if err = db.PingContext(context.TODO()); err != nil {
		panic("Cannot ping db")
	}

	recipeHandler := handler.NewRecipeHandler(repository.NewRecipeRepository(db))
	fileHandler := handler.NewFileHandler()
	ingredientHandler := handler.NewIngredientHandler(repository.NewIngredientRepository(db))

	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/static", internal.StaticFolder)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Ваш фронтенд URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := router.Group("api/v1")

	apiV1.POST("/files/upload", fileHandler.UploadHanlder)
	apiV1.POST("/recipes", recipeHandler.CreateHandler)
	apiV1.GET("/recipes/:id", recipeHandler.GetHandler)
	apiV1.GET("/recipes", recipeHandler.ListHandler)
	apiV1.GET("/ingredients", ingredientHandler.ListHandler)

	apiV1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	fmt.Println("Hello world!")

	router.Run()
}
