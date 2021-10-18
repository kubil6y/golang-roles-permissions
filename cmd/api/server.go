package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func initFlags(cfg *config) {
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 4000, "API Server PORT")

	flag.Parse()
}

func newZapLogger() *zap.SugaredLogger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.Stamp)
	logger, _ := config.Build()
	return logger.Sugar()
}

func (app *application) serve() error {
	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("application start", "server", map[string]string{
		"addr":        srv.Addr,
		"environment": app.config.env,
	})

	if err := srv.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
