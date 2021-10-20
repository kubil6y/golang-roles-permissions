package main

import (
	"fmt"
	"net/http"
)

func (app *application) outOK(message interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ok":   true,
		"data": message,
	}
}

func (app *application) outERR(message interface{}) map[string]interface{} {
	return map[string]interface{}{
		"ok":    false,
		"error": message,
	}
}

func (app *application) logError(r *http.Request, err error) {
	app.logger.Errorw(err.Error(),
		"request_method", r.Method,
		"request_url", r.URL.String(),
	)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	out := app.outERR(message)

	if err := app.writeJSON(w, status, out, nil); err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusInternalServerError, out)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusNotFound, out)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, out)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusTooManyRequests, out)
}

func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusUnauthorized, out)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing token"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusUnauthorized, out)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusUnauthorized, out)
}

func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusForbidden, out)
}

func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	out := app.outERR(message)
	app.errorResponse(w, r, http.StatusForbidden, out)
}
