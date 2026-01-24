package main

import (
	"context"
	"example/techsheets-api/database"
	"example/techsheets-api/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func addIngredient(c *gin.Context) {
	var json pkg.Ingredient
	c.ShouldBindJSON(&json)
	ctx := context.Background()
	gorm.G[pkg.Ingredient](database.DBCon).Create(ctx, &json)
	c.JSON(http.StatusOK, json)
}

func getIngredientById(c *gin.Context) {
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
	c.JSON(http.StatusOK, ingredients)
}

func cellIngredientParser(s *string) {

}

func iterateOverSheets(f *excelize.File) {
	sheets := f.GetSheetList()
	for _, v := range sheets {
		rows, err := f.Rows(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		for rows.Next() {
			row, err := rows.Columns()
			if err != nil {
				fmt.Println(err)
			}
			for key, colCell := range row {
				fmt.Print(key, "\t")
				if ( key == 1){
					
				}
				cellIngredientParser(&colCell)
				fmt.Print(colCell, "\t")
			}
			fmt.Println()
		}
		if err = rows.Close(); err != nil {
			fmt.Println(err)
		}
	}

}

func createIngredientsFromFile(c *gin.Context) {
	file, _ := c.FormFile("file")
	c.SaveUploadedFile(file, "./files/"+file.Filename)
	f, err := excelize.OpenFile("./files/" + file.Filename)
	if err != nil {
		panic(err)
	}

	iterateOverSheets(f)

	c.JSON(http.StatusOK, f.GetSheetList())

}

func main() {
	router := gin.Default()

	router.POST("/ingredient", addIngredient)
	router.GET("/ingredients", getAllIngredients)
	router.POST("/ingredients/upload", createIngredientsFromFile)
	router.GET("/ingredient/:id", getIngredientById)

	db, err := gorm.Open(sqlite.Open("techsheets.db"), &gorm.Config{})
	database.DBCon = db

	if err != nil {
		panic(err)
	}

	// Migrate the schema
	database.DBCon.AutoMigrate(&pkg.Ingredient{})

	log.Println("Table 'Ingredients' created successfully")
	router.Run("localhost:8080")
}
