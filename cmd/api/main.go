package main

import (
	"go.uber.org/zap"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *zap.SugaredLogger
}

func main() {
	var cfg config
	initFlags(&cfg)

	app := &application{
		config: cfg,
		logger: newZapLogger(),
	}

	if err := app.serve(); err != nil {
		app.logger.Fatalf("failed to start %s server", app.config.env)
	}
}
