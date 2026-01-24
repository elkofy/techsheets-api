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

func getIngredientById( c *gin.Context) {
ctx := context.Background()

// Using numeric primary key
ingredient, err := gorm.G[pkg.Ingredient](database.DBCon).Where("id = ?", c.Param("id")).First(ctx)
if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": err})

} else {
	c.JSON(http.StatusOK, ingredient)

 }
}

func removeIngredient(c *gin.Context) {
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
	router.GET("/ingredient/:id", getIngredientById)



	db, err := gorm.Open(sqlite.Open("techsheets.db"), &gorm.Config{})
	database.DBCon = db

	if err != nil {
		panic( err)
	}

	// Migrate the schema
	database.DBCon.AutoMigrate(&pkg.Ingredient{})

	log.Println("Table 'Ingredients' created successfully")
	router.Run("localhost:8080")
}
