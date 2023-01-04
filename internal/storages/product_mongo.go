package storage

import (
	"context"
	"errors"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	db *mongo.Collection
}

func NewProductStorage(db *mongo.Collection) Product {
	return &ProductStorage{db: db}
}

func (p ProductStorage) GetAll(ctx context.Context) ([]domain.Product, error) {
	opts := options.Find()
	opts.SetSort(bson.M{"category.rank": -1})

	cur, err := p.db.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var products []domain.Product
	if err := cur.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (p ProductStorage) GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error) {
	var opts *options.FindOptions
	if sorted {
		opts = options.Find()
		// categories should be sorted with descending order
		opts.SetSort(bson.M{"category.rank": -1})
	}

	cur, err := p.db.Find(ctx, bson.M{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var categories []domain.Category
	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (p ProductStorage) GetCategoryByName(ctx context.Context, categoryName string) (domain.Category, error) {
	res := p.db.FindOne(ctx, bson.M{"category.name": categoryName}, nil)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Category{}, domain.ErrCategoryNotFound
		}
		return domain.Category{}, err
	}
	var category domain.Category
	if err := res.Decode(&category); err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (p ProductStorage) Create(ctx context.Context, product domain.Product) error {
	_, err := p.db.InsertOne(ctx, product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (p ProductStorage) Delete(ctx context.Context, productID string) error {
	result, err := p.db.DeleteOne(ctx, bson.M{"_id": ToObjectID(productID)})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (p ProductStorage) Update(ctx context.Context, productID string, dto dto.UpdateProductDTO) error {
	//TODO implement m:e
	panic("implement me")
}

func (p ProductStorage) Approve(ctx context.Context, productID string) error {
	query := bson.D{{"$set", bson.D{{"approve", true}}}}
	result, err := p.db.UpdateOne(ctx, bson.M{"_id": ToObjectID(productID)}, query)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (p ProductStorage) ChangeImageURL(ctx context.Context, productID string, imageURL string) error {
	query := bson.D{{"$set", bson.D{{"imageURL", imageURL}}}}
	result, err := p.db.UpdateOne(ctx, bson.M{"_id": ToObjectID(productID)}, query)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
