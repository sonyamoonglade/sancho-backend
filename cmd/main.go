package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sonyamoonglade/sancho-backend/database"
	"github.com/sonyamoonglade/sancho-backend/internal/config"
	"github.com/sonyamoonglade/sancho-backend/logger"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	configPath, logsPath, production, strict := readCmdArgs()

	err := logger.NewLogger(logger.Config{
		Out:              []string{logsPath},
		Strict:           strict,
		Production:       production,
		EnableStacktrace: false,
	})
	if err != nil {
		return fmt.Errorf("error instantiating logger: %v", err)
	}

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mongo, err := database.Connect(ctx, cfg.Database.URI, cfg.Database.Name)
	if err != nil {
		return fmt.Errorf("error connecting to mongo: %v", err)
	}
	_ = mongo

	app := fiber.New(fiber.Config{
		Immutable:    true,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return nil
		},
	})

	logger.Get().Info("application is running",
		zap.String("port", cfg.App.Port),
	)
	_ = app

	return app.Listen(":" + cfg.App.Port)
}

func readCmdArgs() (string, string, bool, bool) {

	production := flag.Bool("production", false, "if logger should write to file")
	strict := flag.Bool("strict", false, "if logger should log only warn+ logs")
	logsPath := flag.String("logs-path", "", "where log file is")
	configPath := flag.String("config-path", "", "where config file is")

	flag.Parse()

	// Critical for app if not specified
	if *configPath == "" {
		panic("config path is not provided")
	}

	// Naked return, see return variable names
	return *configPath, *logsPath, *production, *strict
}
