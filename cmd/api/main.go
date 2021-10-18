package main

import (
	"time"

	"github.com/kubil6y/myshop-go/internal/data"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	limiter struct {
		enabled bool
		rps     float64
		burst   int
	}
}

type application struct {
	config config
	logger *zap.SugaredLogger
	models data.Models
}

func main() {
	var cfg config
	initFlags(&cfg)

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Stamp)
	logger, _ := config.Build()
	sugar := logger.Sugar()

	db, err := connectDatabase(cfg)
	if err != nil {
		sugar.Fatal("database connection failed")
	}
	autoMigrate(db)
	sugar.Info("database connection pool established")

	app := &application{
		config: cfg,
		logger: sugar,
		models: data.NewModels(db),
	}

	if err := app.serve(); err != nil {
		app.logger.Fatalf("failed to start %s server", app.config.env)
	}
}
