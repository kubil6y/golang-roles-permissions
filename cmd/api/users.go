package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input registerUserDTO

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var user data.User
	input.populate(&user)

	if err := app.models.Users.Create(&user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, 200, envelope{"user": user}, nil)
	return
}
