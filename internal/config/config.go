package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	service "github.com/sonyamoonglade/sancho-backend/internal/services"
	"github.com/spf13/viper"
)

var (
	ErrConfigNoExist = errors.New("config file doesn't exist")
)

type AppConfig struct {
	Database struct {
		// Connection string
		URI string
		// Name of database
		Name string
	}

	App struct {
		Port string
	}

	Order service.OrderConfig
}

func ReadConfig(path string) (AppConfig, error) {
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return AppConfig{}, ErrConfigNoExist
		}
		return AppConfig{}, err
	}
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		return AppConfig{}, fmt.Errorf("missing MONGO_URI env")
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		return AppConfig{}, fmt.Errorf("missing APP_PORT env")
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		return AppConfig{}, fmt.Errorf("missing DB_NAME env")
	}

	pendingOrderWaitTimeMinutes := viper.GetInt64("order.pending_wait_time")
	if pendingOrderWaitTimeMinutes == 0 {
		return AppConfig{}, fmt.Errorf("missing order.pending_wait_time in config")
	}

	return AppConfig{
		Database: struct {
			URI  string
			Name string
		}{
			URI:  mongoURI,
			Name: dbname,
		},
		App: struct {
			Port string
		}{
			Port: appPort,
		},
		Order: service.OrderConfig{
			PendingOrderWaitTime: time.Duration(pendingOrderWaitTimeMinutes) * time.Minute,
		},
	}, nil
}
