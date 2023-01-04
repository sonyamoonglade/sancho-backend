package domain

import (
	"errors"
)

var (
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrProductNotFound      = errors.New("product not found")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrNoCategories         = errors.New("categories not found")
)

type Product struct {
	ProductID   string   `bson:"_id,omitempty" json:"productId,omitempty"`
	Name        string   `bson:"name" json:"name"`
	TranslateRU string   `bson:"translateRU" json:"translateRu"`
	Description string   `bson:"description" json:"description"`
	ImageURL    string   `bson:"imageURL" json:"imageURL"`
	IsApproved  bool     `bson:"isApproved" json:"isApproved"`
	Price       int64    `bson:"price" json:"price"`
	Category    Category `bson:"category" json:"category"`
	Features    Features `bson:"features" json:"features"`
}

type Category struct {
	CategoryID string
	Rank       int32
	Name       string
}

type Features struct {
	IsLiquid    bool       `bson:"isLiquid" json:"isLiquid"`
	Weight      int32      `bson:"weight" json:"weight,omitempty"`
	Volume      int32      `bson:"volume" json:"volume,omitempty"`
	EnergyValue int32      `bson:"energyValue" json:"energyValue,omitempty"`
	Nutrients   *Nutrients `bson:"nutrients" json:"nutrients,omitempty"`
}

type Nutrients struct {
	Carbs    int32 `bson:"carbs" json:"carbs"`
	Proteins int32 `bson:"proteins" json:"proteins"`
	Fats     int32 `bson:"fats" json:"fats"`
}
