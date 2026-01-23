package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"example/techsheets-api/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type ingredient struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Label     string  `json:"artist"`
	Allergen  string  `json:"allergen"`
	Unit      string  `json:"unit"`
	HT_Price  float64 `json:"ht_price"`
	TTC_price float64 `json:"ttc_price"`
	TVA_rate  float64 `json:"tva_rate"`
	Supplier  string  `json:"suplier"`
}

var ingredients = []ingredient{
	{ID: "1", Name: "Oeuf", Label: "OEUF", Allergen: "oeuf", Unit: "unité", HT_Price: 0.20, TTC_price: 0.30, TVA_rate: 0.055, Supplier: "FERME DE COLLONGES"},
	{ID: "2", Name: "Lait", Label: "LAIT", Allergen: "lactose", Unit: "L", HT_Price: 0.8, TTC_price: 0.90, TVA_rate: 0.055, Supplier: "FERME DE COLLONGES"},
	{ID: "3", Name: "Farine", Label: "FARINET55", Allergen: "gluten", Unit: "KG", HT_Price: 0.35, TTC_price: 0.40, TVA_rate: 0.055, Supplier: "MINOTERIE DES PRèS"},
}

func getIngredients(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, ingredients)
}

func main() {
	router := gin.Default()
	router.GET("/ingredients", getIngredients)
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	ctx := context.Background()

	// Migrate the schema
	db.AutoMigrate(&pkg.IngredientDB{})

	// Create
	// err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})

	// Read
	// product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx) // find product with integer primary key
	// log.Println("Table 'users' created successfully", product)

	//   products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42

	//   // Update - update product's price to 200
	//   products, err = gorm.G[Product](db).Where("id = ?", product.ID).Update(ctx, "Price", 200)
	//   // Update - update multiple fields
	//   products, err = gorm.G[Product](db).Where("id = ?", product.ID).Updates(ctx, Product{Code: "D42", Price: 100})

	//   // Delete - delete product
	//   products, err = gorm.G[Product](db).Where("id = ?", product.ID).Delete(ctx)
	// db, err := sql.Open("sqlite3", "./techsheets.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// sqlStmt := `
	// CREATE TABLE IF NOT EXISTS users (
	//     id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	//     name TEXT
	// );
	// `
	// _, err = db.Exec(sqlStmt)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Println("Table 'users' created successfully")
	router.Run("localhost:8080")
}
