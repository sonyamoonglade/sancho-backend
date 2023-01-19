package tests

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	handler "github.com/sonyamoonglade/sancho-backend/internal/handler"
	"github.com/sonyamoonglade/sancho-backend/internal/handler/middleware"
	service "github.com/sonyamoonglade/sancho-backend/internal/services"
	storage "github.com/sonyamoonglade/sancho-backend/internal/storages"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
	"github.com/sonyamoonglade/sancho-backend/pkg/database"
	"github.com/sonyamoonglade/sancho-backend/pkg/hash"
	"github.com/sonyamoonglade/sancho-backend/pkg/logger"
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

	tokenProvider auth.TokenProvider

	app *fiber.App
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

	app := fiber.New(fiber.Config{
		Immutable:    true,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: handler.HandleError,
	})

	s.app = app
	s.handler.InitAPI(app)
}

func (s *APISuite) TearDownSuite() {
	s.db.Close(context.Background()) //nolint:errcheck
}

func (s *APISuite) initDeps(mongo *database.Mongo) {
	logger.NewLogger(logger.Config{
		Out:              nil,
		Strict:           false,
		Production:       false,
		EnableStacktrace: true,
	})
	logger.Get().Info("Booting e2e test")

	var (
		ttl    = time.Second * 5
		key    = []byte("mama is ok")
		issuer = "localhost"
	)
	tokenProvider, err := auth.NewProvider(ttl, key, issuer)
	if err != nil {
		panic(err)
	}

	storages := storage.NewStorages(mongo)
	services := service.NewServices(service.Deps{
		Storages:      storages,
		TokenProvider: tokenProvider,
		Hasher:        hash.NewSHA1Hasher(),
		TTLStrategy:   ttlStrategy,
	})

	jwtAuth := middleware.NewJWTAuthMiddleware(services.Auth, tokenProvider)
	xReqID := new(middleware.XRequestIDMiddleware)
	middlewares := middleware.NewMiddlewares(jwtAuth, xReqID)

	h := handler.NewHandler(services, middlewares)
	s.handler = h
	s.storages = storages
	s.services = services
	s.tokenProvider = tokenProvider
	s.db = mongo
}

func (s *APISuite) populateDB(ctx context.Context) error {
	_, err := s.db.Collection(storage.CollectionProduct).InsertMany(ctx, products, nil)
	if err != nil {
		return err
	}
	_, err = s.db.Collection(storage.CollectionCategory).InsertMany(ctx, categories, nil)
	if err != nil {
		return err
	}
	_, err = s.db.Collection(storage.CollectionCustomers).InsertOne(ctx, customer, nil)
	if err != nil {
		return err
	}
	return nil
}
