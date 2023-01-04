package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sonyamoonglade/sancho-backend/database"
	handler "github.com/sonyamoonglade/sancho-backend/internal/handler"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/tests/fixtures"
	"github.com/stretchr/testify/suite"
)

var mongoURI, dbName string

func init() {
	mongoURI = os.Getenv("MONGO_URI")
	dbName = os.Getenv("DB_NAME")
}

type APISuite struct {
	suite.Suite

	db       *database.Mongo
	handler  *handler.Handler
	services *service.Services
	storages *storage.Storages
}

func TestAPISuite(t *testing.T) {
	if testing.Short() {
		t.Skipf("skip integration test")
	}

	suite.Run(t, new(APISuite))
}

func (s *APISuite) SetupSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongo, err := database.Connect(ctx, mongoURI, dbName)
	if err != nil {
		s.FailNow("failed to connect to mongodb", err)
		return
	}

	s.initDeps(mongo)
	if err := s.populateDB(ctx); err != nil {
		s.FailNow("failed to populate database", err)
	}
}

func (s *APISuite) initDeps(mongo *database.Mongo) {
	storages := storage.NewStorages(s.db)
	services := service.NewServices(service.Deps{Storages: storages})
	h := handler.NewHandler(services)

	s.handler = h
	s.storages = storages
	s.services = services
	s.db = mongo
}

func (s *APISuite) populateDB(ctx context.Context) error {
	products := fixtures.GetProducts(10)
	_, err := s.db.Collection(storage.CollectionProduct).InsertMany(ctx, products, nil)
	if err != nil {
		return err
	}
	return nil
}
