package storage

import (
	"errors"

	"github.com/sonyamoonglade/sancho-backend/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionProduct string = "product"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Storages struct {
	Product Product
}

func NewStorages(db *database.Mongo) *Storages {
	return &Storages{
		Product: NewProductStorage(db.Collection(CollectionProduct)),
	}
}

func ToObjectID(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}
