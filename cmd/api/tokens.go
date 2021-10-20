package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input createAuthenticationTokenDto
	if err := app.readJSON(w, r, &input); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	matches, err := user.ComparePassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !matches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeAuthentication)

	e := envelope{"authentication_token": map[string]interface{}{
		"token":  token.Plaintext,
		"expiry": token.Expiry,
	}}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
