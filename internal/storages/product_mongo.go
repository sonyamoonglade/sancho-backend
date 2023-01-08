package storage

import (
	"context"
	"errors"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStorage struct {
	products   *mongo.Collection
	categories *mongo.Collection
}

func NewProductStorage(products *mongo.Collection, categories *mongo.Collection) Product {
	return &ProductStorage{
		products:   products,
		categories: categories,
	}
}

func (p ProductStorage) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	result := p.products.FindOne(ctx, bson.M{"_id": ToObjectID(productID)}, nil)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Product{}, ErrNotFound
		}

		return domain.Product{}, err
	}
	var product domain.Product
	if err := result.Decode(&product); err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p ProductStorage) GetAll(ctx context.Context) ([]domain.Product, error) {
	opts := options.Find()
	opts.SetSort(bson.M{"category.rank": -1})

	cur, err := p.products.Find(ctx, bson.M{}, opts)
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
		opts.SetSort(bson.M{"rank": -1})
	}

	cur, err := p.categories.Find(ctx, bson.M{}, opts)
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
	res := p.categories.FindOne(ctx, bson.M{"name": categoryName}, nil)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Category{}, ErrNotFound
		}
		return domain.Category{}, err
	}
	var category domain.Category
	if err := res.Decode(&category); err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (p ProductStorage) Create(ctx context.Context, product domain.Product) (primitive.ObjectID, error) {
	r, err := p.products.InsertOne(ctx, product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, ErrAlreadyExists
		}
		return primitive.ObjectID{}, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func (p ProductStorage) Delete(ctx context.Context, productID string) error {
	result, err := p.products.DeleteOne(ctx, bson.M{"_id": ToObjectID(productID)})
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
	query := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{
				Key:   "isApproved",
				Value: true,
			},
		}},
	}
	result, err := p.products.UpdateOne(ctx, bson.M{"_id": ToObjectID(productID)}, query)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (p ProductStorage) Disapprove(ctx context.Context, productID string) error {
	query := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{
				Key:   "isApproved",
				Value: false,
			},
		}},
	}
	result, err := p.products.UpdateOne(ctx, bson.M{"_id": ToObjectID(productID)}, query)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (p ProductStorage) ChangeImageURL(ctx context.Context, productID string, imageURL string) error {
	query := bson.D{
		bson.E{Key: "$set", Value: bson.D{
			bson.E{
				Key:   "imageURL",
				Value: imageURL,
			},
		}},
	}
	result, err := p.products.UpdateOne(ctx, bson.M{"_id": ToObjectID(productID)}, query)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
