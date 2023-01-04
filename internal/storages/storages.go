package storage

import "github.com/sonyamoonglade/sancho-backend/database"

const (
	CollectionProduct string = "product"
)

type Storages struct {
	Product Product
}

func NewStorages(db *database.Mongo) *Storages {
	return &Storages{
		Product: NewProductStorage(db.Collection(CollectionProduct)),
	}
}
