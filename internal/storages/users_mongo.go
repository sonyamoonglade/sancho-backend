package storage

import (
	"context"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongo struct {
	users *mongo.Collection
}

func NewUserMongo(users *mongo.Collection) User {
	return userMongo{
		users: users,
	}
}

func (u userMongo) GetAdminByRefreshToken(ctx context.Context, token string) (domain.Admin, error) {
	//TODO implement me
	panic("implement me")
}
