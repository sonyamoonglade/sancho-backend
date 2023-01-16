package storage

import (
	"context"
	"errors"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userStorage struct {
	customers        *mongo.Collection
	adminsAndWorkers *mongo.Collection
}

func NewUserStorage(customers *mongo.Collection, adminsAndWorkers *mongo.Collection) User {
	return userStorage{
		customers:        customers,
		adminsAndWorkers: adminsAndWorkers,
	}
}

func (u userStorage) GetAdminByLogin(ctx context.Context, login string) (domain.Admin, error) {
	//TODO implement me
	panic("implement me")
}

func (u userStorage) GetAdminByRefreshToken(ctx context.Context, adminID, token string) (domain.Admin, error) {

	query := bson.D{bson.E{
		Key:   "_id",
		Value: adminID,
	}, bson.E{
		Key:   "session.refreshToken",
		Value: token,
	}, bson.E{
		Key:   "session.expiresAt",
		Value: bson.M{"$gt": time.Now().UTC()},
	}}

	result := u.adminsAndWorkers.FindOne(ctx, query)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.Admin{}, domain.ErrAdminNotFound
		}
		return domain.Admin{}, err
	}

	var admin domain.Admin
	if err := result.Decode(&admin); err != nil {
		return domain.Admin{}, err
	}

	return admin, nil
}

func (u userStorage) SaveAdmin(ctx context.Context, admin domain.Admin) (string, error) {
	res, err := u.adminsAndWorkers.InsertOne(ctx, admin)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", domain.ErrAdminAlreadyExists
		}
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (u userStorage) SaveCustomer(ctx context.Context, customer domain.Customer) error {
	//TODO implement me
	panic("implement me")
}

func (u userStorage) SaveWorker(ctx context.Context, worker domain.Worker) error {
	//TODO implement me
	panic("implement me")
}

func (u userStorage) SaveSession(ctx context.Context, dto dto.SaveSessionDTO) error {

	var collection = u.customers
	if dto.Role == domain.RoleAdmin || dto.Role == domain.RoleWorker {
		collection = u.adminsAndWorkers
	}

	opts := options.Update()
	opts.SetUpsert(true)

	updateQuery := bson.D{bson.E{
		Key:   "$set",
		Value: dto.Session,
	}}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": ToObjectID(dto.UserID)}, updateQuery, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}
