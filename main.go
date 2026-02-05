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
	"strconv"
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

func iterateOverSheets(f *excelize.File) {
	sheets := f.GetSheetList()
	// tables, err := f.GetTables("Feuille1")
	// if err != nil {
	// 	panic(err)
	// }
	// table := tables[0]

	for _, v := range sheets {
		rows, err := f.Rows(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		var columnOptions excelize.Options
		columnOptions.RawCellValue = false
		for rows.Next() {
			row, err := rows.Columns(columnOptions)

			if row[0] == "Ref" {
				rows.Next()
			}
			fmt.Println(row)

			ctx := context.Background()
			var newIngredient pkg.Ingredient
			newIngredientPrice, priceErr := strconv.ParseFloat(row[3], 64)
			if priceErr != nil {
				newIngredientPrice = 0
			}
			// newIngredientEnergy, energyErr := strconv.ParseFloat(row[5], 64)
			// if energyErr != nil {
			// 	newIngredientEnergy = 0
			// }
			newIngredient.Ref = row[0]
			newIngredient.Name = row[1]
			newIngredient.Unit = row[2]

			newIngredient.Unit_Price = newIngredientPrice
			newIngredient.Category = row[4]
			// newIngredient.Energy= newIngredientEnergy
			// newIngredient.Supplier = row[6]
			// newIngredient.Allergen = row[7]
			// newIngredient.Details = row[8]

			gorm.G[pkg.Ingredient](database.DBCon).Create(ctx, &newIngredient)
			if err != nil {
				fmt.Println(err)
			}
			for _, colCell := range row {
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

func deleteIngredientById(c *gin.Context) {
	ctx := context.Background()
	rowsAffected, err := gorm.G[pkg.Ingredient](database.DBCon).Where("id = ?", c.Param("id")).Delete(ctx)
	if err == nil {
		c.JSON(http.StatusNoContent, gin.H{
			"code":    http.StatusNoContent,
			"message": "The Ingredient with id :" + c.Param("id") + " was removed",
		})
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotModified, gin.H{
			"code":    http.StatusNotModified,
			"message": "The Ingredient with id :" + c.Param("id") + " was not removed",
		})
	}
}

func updateIngredientById(c *gin.Context) {
	ctx := context.Background()
	var json pkg.Ingredient
	c.ShouldBindJSON(&json)
	ingredient, err := gorm.G[pkg.Ingredient](database.DBCon).Where("id = ?", c.Param("id")).Updates(ctx, json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, ingredient)
	}
}
// Source - https://stackoverflow.com/a/63811206
// Posted by Rahul S, modified by community. See post 'Timeline' for change history
// Retrieved 2026-02-05, License - CC BY-SA 4.0

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {

        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Credentials", "true")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}


func main() {
	router := gin.Default()
	
	// ✅ Ajouter le middleware CORS EN PREMIER
	router.Use(CORSMiddleware())
	
	// ✅ Ensuite enregistrer les routes
	router.POST("/ingredient", addIngredient)
	router.GET("/ingredients", getAllIngredients)
	router.POST("/ingredients/upload", createIngredientsFromFile)
	router.GET("/ingredient/:id", getIngredientById)
	router.DELETE("/ingredient/:id", deleteIngredientById)
	router.PATCH("/ingredient/:id", updateIngredientById)
	
	db, err := gorm.Open(sqlite.Open("techsheets.db"), &gorm.Config{})
	database.DBCon = db
	if err != nil {
		panic(err)
	}
	
	database.DBCon.AutoMigrate(&pkg.Ingredient{})
	log.Println("Table 'Ingredients' created successfully")
	
	router.Run("localhost:8080")
}
