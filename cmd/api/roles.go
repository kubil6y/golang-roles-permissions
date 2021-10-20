package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

func (app *application) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	var input createRoleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var role data.Role
	input.populate(&role)
	if err := app.models.Roles.Insert(&role); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("name", "a role with that name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"role": role}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllRolesHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {}

func (app *application) updateRolesHandler(w http.ResponseWriter, r *http.Request) {}
