package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func initFlags(cfg *config) {
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 4000, "API Server PORT")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MYSHOP_DB_DSN"), "PostgreSQL DSN")

	flag.Parse()
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
