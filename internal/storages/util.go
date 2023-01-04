package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

func ToObjectID(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}
