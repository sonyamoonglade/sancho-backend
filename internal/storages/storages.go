package storage

import (
	"strings"

	"github.com/sonyamoonglade/sancho-backend/pkg/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionOrders           = "orders"
	CollectionProduct          = "products"
	CollectionCategory         = "categories"
	CollectionCustomers        = "customers"
	CollectionAdminsAndWorkers = "adminsAndWorkers"
)

type Storages struct {
	Product Product
	User    User
	Order   Order
}

func NewStorages(db *database.Mongo) *Storages {
	return &Storages{
		Product: NewProductStorage(db.Collection(CollectionProduct), db.Collection(CollectionCategory)),
		User:    NewUserStorage(db.Collection(CollectionCustomers), db.Collection(CollectionAdminsAndWorkers)),
		Order:   NewOrderStorage(db.Collection(CollectionOrders)),
	}
}

func ToObjectID(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}

func GetFieldAndValueFromDuplicateError(err error) (field string, value string) {
	var (
		msg           = err.Error()
		split         = strings.Split(msg, "{")
		textToProcess = strings.TrimSpace(split[1])
		fieldDone     bool
		skip          bool
	)
	for _, ch := range strings.Split(textToProcess, "") {
		// When meet ':' skip next 2 chars (space and double quote)
		if ch == ":" {
			fieldDone = true
			skip = true

			// 1st skip
			continue
		}
		if skip {
			skip = false
			// 2nd skip
			continue
		}
		if !fieldDone {
			field += ch
			continue
		}
		if ch != `\` && ch != `"` {
			if ch == "}" {
				value = value[:len(value)-1]
				break
			}
			value += ch
		}
	}
	return field, value
}
