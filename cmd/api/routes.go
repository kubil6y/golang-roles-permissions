package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.HandlerFunc(http.MethodGet, "/v1/api/healthcheck", app.healthCheckHandler)

	return app.recoverPanic(app.rateLimit(router))
}
