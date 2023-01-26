package storage

import (
	"context"
	"errors"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/appErrors"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type orderStorage struct {
	orders *mongo.Collection
}

func NewOrderStorage(orders *mongo.Collection) Order {
	return &orderStorage{orders: orders}
}

func (o orderStorage) GetOrderByID(ctx context.Context, orderID string) (domain.Order, error) {
	result := o.orders.FindOne(ctx, bson.M{"_id": ToObjectID(orderID)}, nil)
	if err := result.Err(); err != nil {
		return domain.Order{}, err
	}
	var order domain.Order
	if err := result.Decode(&order); err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (o orderStorage) GetLastOrderByCustomerID(ctx context.Context, customerID string) (domain.Order, error) {
	opts := options.FindOne()
	opts.SetSort(bson.M{"createdAt": -1})

	query := bson.D{bson.E{
		Key:   "customerId",
		Value: ToObjectID(customerID),
	}}

	result := o.orders.FindOne(ctx, query)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return domain.Order{}, err
	}

	var order domain.Order
	if err := result.Decode(&order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (o orderStorage) GetOrderByNanoIDAt(ctx context.Context, nanoID string, from, to time.Time) (domain.Order, error) {
	var (
		isoTo   = to.Format(time.RFC3339)
		isoFrom = from.Format(time.RFC3339)
	)

	query := bson.D{bson.E{
		Key: "createdAt",
		Value: bson.M{
			"$gte": isoFrom,
			"$lte": isoTo,
		}},
		bson.E{
			Key:   "nanoId",
			Value: nanoID,
		},
	}

	result := o.orders.FindOne(ctx, query)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Order{}, domain.ErrOrderNotFound
		}
		return domain.Order{}, err
	}

	var order domain.Order
	if err := result.Decode(&order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (o orderStorage) SaveOrder(ctx context.Context, order domain.Order) (primitive.ObjectID, error) {
	result, err := o.orders.InsertOne(ctx, order)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			field, value := GetFieldAndValueFromDuplicateError(err)
			return primitive.ObjectID{}, appErrors.NewDuplicateError("order", field, value)
		}
		return primitive.ObjectID{}, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
