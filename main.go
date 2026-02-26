package main

import (
	"context"
	_ "example/techsheets-api/docs"
	"example/techsheets-api/database"
	"example/techsheets-api/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// @Summary     Create ingredient
// @Tags        ingredients
// @Accept      json
// @Produce     json
// @Param       ingredient body     pkg.Ingredient true "Ingredient to create"
// @Success     200        {object} pkg.Ingredient
// @Router      /ingredient [post]
func addIngredient(c *gin.Context) {
	var json pkg.Ingredient
	c.ShouldBindJSON(&json)
	ctx := context.Background()
	gorm.G[pkg.Ingredient](database.DBCon).Create(ctx, &json)
	c.JSON(http.StatusOK, json)
}

// @Summary     Get ingredient by ID
// @Tags        ingredients
// @Produce     json
// @Param       id  path     int true "Ingredient ID"
// @Success     200 {object} pkg.Ingredient
// @Failure     400 {object} map[string]interface{}
// @Router      /ingredient/{id} [get]
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

// @Summary     List all ingredients
// @Tags        ingredients
// @Produce     json
// @Success     200 {array} pkg.Ingredient
// @Router      /ingredients [get]
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

// @Summary     Bulk import ingredients from Excel file
// @Tags        ingredients
// @Accept      multipart/form-data
// @Produce     json
// @Param       file formData file   true "Excel file (.xlsx)"
// @Success     200  {array}  string
// @Router      /ingredients/upload [post]
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

// @Summary     Delete ingredient by ID
// @Tags        ingredients
// @Param       id  path int true "Ingredient ID"
// @Success     204
// @Failure     304
// @Router      /ingredient/{id} [delete]
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

// @Summary     Update ingredient by ID
// @Tags        ingredients
// @Accept      json
// @Produce     json
// @Param       id         path     int            true "Ingredient ID"
// @Param       ingredient body     pkg.Ingredient true "Updated ingredient"
// @Success     200        {object} pkg.Ingredient
// @Failure     400        {object} map[string]interface{}
// @Router      /ingredient/{id} [patch]
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

// Recipe handlers

// @Summary     Create recipe
// @Tags        recipes
// @Accept      json
// @Produce     json
// @Param       recipe body     pkg.Recipe true "Recipe to create"
// @Success     200    {object} pkg.Recipe
// @Failure     400    {object} map[string]interface{}
// @Router      /recipe [post]
func addRecipe(c *gin.Context) {
	var recipe pkg.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DBCon.Create(&recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// @Summary     List all recipes
// @Tags        recipes
// @Produce     json
// @Success     200 {array} pkg.Recipe
// @Router      /recipes [get]
func getAllRecipes(c *gin.Context) {
	var recipes []pkg.Recipe
	database.DBCon.Find(&recipes)
	c.JSON(http.StatusOK, recipes)
}

// @Summary     Get recipe by ID
// @Tags        recipes
// @Produce     json
// @Param       id  path     int true "Recipe ID"
// @Success     200 {object} pkg.Recipe
// @Failure     404 {object} map[string]interface{}
// @Router      /recipe/{id} [get]
func getRecipeById(c *gin.Context) {
	var recipe pkg.Recipe
	result := database.DBCon.
		Preload("RecipeIngredients.Ingredient").
		Preload("Steps").
		First(&recipe, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// @Summary     Update recipe by ID
// @Tags        recipes
// @Accept      json
// @Produce     json
// @Param       id     path     int        true "Recipe ID"
// @Param       recipe body     pkg.Recipe true "Updated recipe"
// @Success     200    {object} pkg.Recipe
// @Failure     400    {object} map[string]interface{}
// @Router      /recipe/{id} [patch]
func updateRecipeById(c *gin.Context) {
	var recipe pkg.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DBCon.Model(&pkg.Recipe{}).Where("id = ?", c.Param("id")).Updates(recipe).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// @Summary     Delete recipe by ID
// @Tags        recipes
// @Param       id  path int true "Recipe ID"
// @Success     204
// @Failure     400 {object} map[string]interface{}
// @Router      /recipe/{id} [delete]
func deleteRecipeById(c *gin.Context) {
	if err := database.DBCon.Delete(&pkg.Recipe{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// TechnicalSheet handlers

// @Summary     Create technical sheet
// @Tags        technical-sheets
// @Accept      json
// @Produce     json
// @Param       sheet body     pkg.TechnicalSheet true "Technical sheet to create"
// @Success     200   {object} pkg.TechnicalSheet
// @Failure     400   {object} map[string]interface{}
// @Router      /technicalsheet [post]
func addTechnicalSheet(c *gin.Context) {
	var sheet pkg.TechnicalSheet
	if err := c.ShouldBindJSON(&sheet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DBCon.Create(&sheet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sheet)
}

// @Summary     List all technical sheets
// @Tags        technical-sheets
// @Produce     json
// @Success     200 {array} pkg.TechnicalSheet
// @Router      /technicalsheets [get]
func getAllTechnicalSheets(c *gin.Context) {
	var sheets []pkg.TechnicalSheet
	database.DBCon.Find(&sheets)
	c.JSON(http.StatusOK, sheets)
}

// @Summary     Get technical sheet by ID
// @Tags        technical-sheets
// @Produce     json
// @Param       id  path     int true "Technical sheet ID"
// @Success     200 {object} pkg.TechnicalSheet
// @Failure     404 {object} map[string]interface{}
// @Router      /technicalsheet/{id} [get]
func getTechnicalSheetById(c *gin.Context) {
	var sheet pkg.TechnicalSheet
	result := database.DBCon.
		Preload("Mold").
		Preload("Recipes.RecipeIngredients.Ingredient").
		Preload("Recipes.Steps").
		Preload("FinishingSteps").
		Preload("AdditionalRecipeSteps").
		First(&sheet, c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, sheet)
}

// @Summary     Update technical sheet by ID
// @Tags        technical-sheets
// @Accept      json
// @Produce     json
// @Param       id    path     int                true "Technical sheet ID"
// @Param       sheet body     pkg.TechnicalSheet true "Updated technical sheet"
// @Success     200   {object} pkg.TechnicalSheet
// @Failure     400   {object} map[string]interface{}
// @Router      /technicalsheet/{id} [patch]
func updateTechnicalSheetById(c *gin.Context) {
	var sheet pkg.TechnicalSheet
	if err := c.ShouldBindJSON(&sheet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DBCon.Model(&pkg.TechnicalSheet{}).Where("id = ?", c.Param("id")).Updates(sheet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sheet)
}

// @Summary     Delete technical sheet by ID
// @Tags        technical-sheets
// @Param       id  path int true "Technical sheet ID"
// @Success     204
// @Failure     400 {object} map[string]interface{}
// @Router      /technicalsheet/{id} [delete]
func deleteTechnicalSheetById(c *gin.Context) {
	if err := database.DBCon.Delete(&pkg.TechnicalSheet{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
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

// @title           TechSheets API
// @version         1.0
// @description     REST API for managing ingredients, recipes and technical sheets.
// @host            localhost:8080
// @BasePath        /
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

	router.POST("/recipe", addRecipe)
	router.GET("/recipes", getAllRecipes)
	router.GET("/recipe/:id", getRecipeById)
	router.PATCH("/recipe/:id", updateRecipeById)
	router.DELETE("/recipe/:id", deleteRecipeById)

	router.POST("/technicalsheet", addTechnicalSheet)
	router.GET("/technicalsheets", getAllTechnicalSheets)
	router.GET("/technicalsheet/:id", getTechnicalSheetById)
	router.PATCH("/technicalsheet/:id", updateTechnicalSheetById)
	router.DELETE("/technicalsheet/:id", deleteTechnicalSheetById)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	db, err := gorm.Open(sqlite.Open("techsheets.db"), &gorm.Config{})
	database.DBCon = db
	if err != nil {
		panic(err)
	}

	database.DBCon.AutoMigrate(
		&pkg.Ingredient{},
		&pkg.Mold{},
		&pkg.Recipe{},
		&pkg.RecipeIngredient{},
		&pkg.Step{},
		&pkg.AdditionalRecipeSteps{},
		&pkg.TechnicalSheet{},
	)
	log.Println("AutoMigrate completed for all tables")

	router.Run("localhost:8080")
}
