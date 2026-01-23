package main

import (
	"context"
	"example/techsheets-api/database"
	"example/techsheets-api/pkg"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"fmt"
)


func addIngredient(c *gin.Context) {
	var json pkg.Ingredient
	c.ShouldBindJSON(&json)
	ctx := context.Background()
	gorm.G[pkg.Ingredient](database.DBCon).Create(ctx, &json)
	c.JSON(http.StatusOK, json)
}

func getAllIngredients(c *gin.Context) {
		var ingredients []pkg.Ingredient
		result := database.DBCon.Find(&ingredients)
		
    fmt.Println("Hello, World!", result)
	c.JSON(http.StatusOK, ingredients	)
}

func main() {
	router := gin.Default()

	router.POST("/ingredient", addIngredient)
	router.GET("/ingredients", getAllIngredients)


	db, err := gorm.Open(sqlite.Open("techsheets.db"), &gorm.Config{})
	database.DBCon = db

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	database.DBCon.AutoMigrate(&pkg.Ingredient{})

	log.Println("Table 'users' created successfully")
	router.Run("localhost:8080")
}
