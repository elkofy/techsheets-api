package pkg

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DOMAIN
// type Ingredient struct {
// 	UUID          uuid.UUID
// 	Name          string
// 	Label         string
// 	Allergen      string
// 	Unit          string
// 	HT_Price      float64
// 	TTC_price     float64
// 	TVA_rate      float64
// 	Supplier      string
// }

type Ingredient struct {
	gorm.Model
	Name       string  `form:"name" json:"name" `
	Ref        string  `form:"ref" json:"ref"`
	Allergen   string  `form:"allergen" json:"allergen"`
	Unit       string  `form:"unit" json:"unit"`
	Unit_Price float64 `form:"unitPrice" json:"unitPrice"`
	Category   string  `form:"category" json:"category"`
	Details    string  `form:"details" json:"details"`
	Supplier   string  `form:"supplier" json:"supplier"`
	Energy     float64 `form:"energy" json:"energy"`
}

type RecipeIngredient struct {
	gorm.Model
	RecipeID     uint       `json:"recipeId"`
	IngredientID uint       `json:"ingredientId"`
	Quantity     float64    `json:"quantity"`
	Ingredient   Ingredient `json:"ingredient" gorm:"foreignKey:IngredientID"`
}

type Verb string

const (
	VerbChop Verb = "Couper"
	VerbMix  Verb = "Melanger"
	VerbBake Verb = "Cuire"
	VerbBoil Verb = "Amalgamer"
	VerbFry  Verb = "Monter"
)

type Step struct {
	gorm.Model
	Verb             Verb   `json:"verb"`
	ActionDetail     string `json:"actionDetail"`
	VideoUrl         string `json:"videoUrl"`
	RecipeID         *uint  `json:"recipeId,omitempty"`
	TechnicalSheetID *uint  `json:"technicalSheetId,omitempty"`
}

type Recipe struct {
	gorm.Model
	Name              string             `json:"name"`
	RecipeIngredients []RecipeIngredient `json:"recipeIngredients" gorm:"foreignKey:RecipeID"`
	Steps             []Step             `json:"steps" gorm:"foreignKey:RecipeID"`
	Timings           Timings            `json:"timings" gorm:"serializer:json"`
	Equipements       []string           `json:"equipements" gorm:"serializer:json"`
}

type Mold struct {
	gorm.Model
	Name     string  `json:"name"`
	Shape    string  `json:"shape"`
	Capacity float64 `json:"capacity"`
}

type AdditionalRecipeSteps struct {
	gorm.Model
	RecipeID         uint `json:"recipeId"`
	TechnicalSheetID uint `json:"technicalSheetId"`
	StepID           uint `json:"stepId"`
}

type TechnicalSheet struct {
	gorm.Model
	Name                  string                  `json:"name"`
	Description           string                  `json:"description"`
	MoldID                *uint                   `json:"moldId,omitempty"`
	Mold                  Mold                    `json:"mold" gorm:"foreignKey:MoldID"`
	Recipes               []Recipe                `json:"recipes" gorm:"many2many:technical_sheet_recipes;"`
	Conservation          string                  `json:"conservation"`
	ImageUrl              string                  `json:"imageUrl"`
	FinishingSteps        []Step                  `json:"finishingSteps" gorm:"foreignKey:TechnicalSheetID"`
	AdditionalRecipeSteps []AdditionalRecipeSteps `json:"additionalRecipeSteps" gorm:"foreignKey:TechnicalSheetID"`
	Yield                 Yield                   `json:"yield" gorm:"serializer:json"`
}

type Type string

const (
	TypeSlice    Type = "Part"
	TypeRamequin Type = "Ramequin"
)

type Portion struct {
	Type     Type  `json:"type"`
	Quantity int32 `json:"quantity"`
}

type Yield struct {
	Servings int32   `json:"servings"`
	Portion  Portion `json:"portion"`
}

type Timings struct {
	Preparation Timing     `json:"preparation"`
	Cooking     Timing     `json:"cooking"`
	Baking      BakeTiming `json:"baking"`
	Resting     Timing     `json:"resting"`
	TotalTime   int        `json:"totalTime"`
}

type Timing struct {
	Duration    int    `json:"duration"`
	Unit        string `json:"unit"`
	Description string `json:"description,omitempty"`
}

type BakeTiming struct {
	Duration    int         `json:"duration"`
	Unit        string      `json:"unit"`
	Temperature Temperature `json:"temperature"`
}

type Temperature struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

// DTO

type IngredientDTO struct {
	UUID      uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Label     string    `json:"artist"`
	Allergen  string    `json:"allergen"`
	Unit      string    `json:"unit"`
	HT_Price  float64   `json:"ht_price"`
	TTC_price float64   `json:"ttc_price"`
	TVA_rate  float64   `json:"tva_rate"`
	Supplier  string    `json:"suplier"`
}
