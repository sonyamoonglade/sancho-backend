package domain

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrProductAlreadyExists      = errors.New("product already exists")
	ErrProductAlreadyApproved    = errors.New("product is already approved")
	ErrProductAlreadyDisapproved = errors.New("product is already disapproved")
	ErrProductNotFound           = errors.New("product not found")
	ErrCategoryNotFound          = errors.New("category not found")
	ErrNoCategories              = errors.New("categories not found")
)

type Product struct {
	ProductID   primitive.ObjectID `bson:"_id,omitempty" json:"productId,omitempty"`
	Name        string             `bson:"name" json:"name"`
	TranslateRU string             `bson:"translateRu" json:"translateRu"`
	Description string             `bson:"description" json:"description"`
	ImageURL    *string            `bson:"imageUrl" json:"imageUrl"`
	IsApproved  bool               `bson:"isApproved" json:"isApproved"`
	Price       int64              `bson:"price" json:"price"`
	Category    Category           `bson:"category" json:"category"`
	Features    Features           `bson:"features" json:"features"`
}

type Category struct {
	CategoryID primitive.ObjectID `bson:"_id" json:"categoryId"`
	Rank       int32              `bson:"rank" json:"rank"`
	Name       string             `bson:"name" json:"name"`
}

type Features struct {
	IsLiquid    bool       `bson:"isLiquid" json:"isLiquid"`
	Weight      int32      `bson:"weight" json:"weight,omitempty"`
	Volume      int32      `bson:"volume" json:"volume,omitempty"`
	EnergyValue int32      `bson:"energyValue" json:"energyValue,omitempty"`
	Nutrients   *Nutrients `bson:"nutrients,omitempty" json:"nutrients,omitempty"`
}

type Nutrients struct {
	Carbs    int32 `bson:"carbs,omitempty" json:"carbs"`
	Proteins int32 `bson:"proteins,omitempty" json:"proteins"`
	Fats     int32 `bson:"fats,omitempty" json:"fats"`
}
