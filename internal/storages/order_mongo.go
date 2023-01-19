package storage

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type orderStorage struct {
	orders *mongo.Collection
}

func NewOrderStorage(orders *mongo.Collection) Order {
	return &orderStorage{orders: orders}
}

func (o orderStorage) SaveOrder(ctx context.Context, order domain.Order) (string, error) {
	//TODO implement me
	panic("implement me")
}
