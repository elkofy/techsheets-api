package pkg

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
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
	Details  string `form:"details" json:"details"`
	Supplier   string  `form:"supplier" json:"supplier"`
	Energy   float64  `form:"energy" json:"energy"`
}

type RecipeIngredient struct {
	UUID           uuid.UUID
	RecipeUUID     uuid.UUID
	IngredientUUID uuid.UUID
	Quantity       float64
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
	UUID         uuid.UUID
	Verb         Verb
	actionDetail string
	videoUrl     string
}

type RecipeTime struct {
	weighing      time.Time
	prepping      time.Time
	making        time.Time
	conditionning time.Time
}

type Recipe struct {
	UUID            uuid.UUID
	Name            string
	IngredientUUIDs []uuid.UUID
	Steps           []Step
	Timings         Timings
	Equipements     []string
	createdAt       time.Time
	lastUpdatedAt   time.Time
}

type Mold struct {
	UUID     uuid.UUID
	Name     string
	Shape    string
	Capacity float64
}

type AdditionalRecipeSteps struct {
	UUID               uuid.UUID
	RecipeUUID         uuid.UUID
	TechnicalSheetUUID uuid.UUID
	StepUUID           uuid.UUID
}

type TechnicalSheet struct {
	UUID                  uuid.UUID
	Name                  string
	Description           string
	Mold                  Mold
	Recipes               []Recipe
	Conservation          string
	ImageUrl              string
	FinishingSteps        []Step
	AdditionalRecipeSteps []AdditionalRecipeSteps
	Yield                 Yield
}

type Type string

const (
	TypeSlice    Type = "Part"
	TypeRamequin Type = "Ramequin"
)

type Portion struct {
	Type     Type
	quantity int32
}
type Yield struct {
	Servings int32
	Portion  Portion
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

//DTO

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
