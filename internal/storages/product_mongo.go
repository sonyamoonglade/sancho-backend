package storage

import (
	"context"
	"errors"

	"github.com/sonyamoonglade/sancho-backend/internal/appErrors"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productStorage struct {
	products   *mongo.Collection
	categories *mongo.Collection
}

func NewProductStorage(products *mongo.Collection, categories *mongo.Collection) Product {
	return &productStorage{
		products:   products,
		categories: categories,
	}
}

func (p productStorage) GetByID(ctx context.Context, productID string) (domain.Product, error) {
	result := p.products.FindOne(ctx, bson.M{"_id": ToObjectID(productID)}, nil)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Product{}, domain.ErrProductNotFound
		}

		return domain.Product{}, err
	}
	var product domain.Product
	if err := result.Decode(&product); err != nil {
		return domain.Product{}, err
	}
	return product, nil
}

func (p productStorage) GetAll(ctx context.Context) ([]domain.Product, error) {
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
func (p productStorage) GetByIDs(ctx context.Context, ids []string) ([]domain.Product, error) {
	objectIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		objectIDs = append(objectIDs, ToObjectID(id))
	}

	query := bson.D{bson.E{
		Key:   "_id",
		Value: bson.M{"$in": objectIDs},
	}}

	cursor, err := p.products.Find(ctx, query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoProducts
		}
		return nil, err
	}

	products := make([]domain.Product, 0, len(ids))
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (p productStorage) Save(ctx context.Context, product domain.Product) (primitive.ObjectID, error) {
	r, err := p.products.InsertOne(ctx, product)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return primitive.ObjectID{}, domain.ErrProductAlreadyExists
		}
		return primitive.ObjectID{}, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func (p productStorage) Delete(ctx context.Context, productID string) error {
	result, err := p.products.DeleteOne(ctx, bson.M{"_id": ToObjectID(productID)})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (p productStorage) Update(ctx context.Context, dto dto.UpdateProductDTO) error {
	updateQuery := bson.M{}
	if dto.Price != nil {
		updateQuery["price"] = *dto.Price
	}
	if dto.Name != nil {
		updateQuery["name"] = *dto.Name
	}
	if dto.Description != nil {
		updateQuery["description"] = *dto.Description
	}
	if dto.TranslateRU != nil {
		updateQuery["translateRu"] = *dto.TranslateRU
	}
	if dto.ImageURL != nil {
		updateQuery["imageUrl"] = *dto.ImageURL
	}
	setQuery := bson.D{
		bson.E{
			Key:   "$set",
			Value: updateQuery,
		},
	}
	result, err := p.products.UpdateOne(ctx, bson.M{"_id": ToObjectID(dto.ProductID)}, setQuery, nil)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			field, value := GetFieldAndValueFromDuplicateError(err)
			return appErrors.NewDuplicateError("product", field, value)
		}
		return err
	}
	if result.MatchedCount == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (p productStorage) Approve(ctx context.Context, productID string) error {
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
		return domain.ErrProductNotFound
	}
	return nil
}

func (p productStorage) Disapprove(ctx context.Context, productID string) error {
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
		return domain.ErrProductNotFound
	}
	return nil
}

func (p productStorage) GetAllCategories(ctx context.Context, sorted bool) ([]domain.Category, error) {
	var opts *options.FindOptions
	if sorted {
		opts = options.Find()
		// categories should be sorted with descending order
		opts.SetSort(bson.M{"rank": -1})
	}

	cur, err := p.categories.Find(ctx, bson.M{}, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrNoCategories
		}
		return nil, err
	}

	var categories []domain.Category
	if err := cur.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}

func (p productStorage) GetCategoryByName(ctx context.Context, categoryName string) (domain.Category, error) {
	res := p.categories.FindOne(ctx, bson.M{"name": categoryName}, nil)
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
