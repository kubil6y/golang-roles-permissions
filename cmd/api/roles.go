package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

func (app *application) createRoleHandler(w http.ResponseWriter, r *http.Request) {
	var input roleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// check if permissions exists
	permissions := make([]data.Permission, 0)
	for _, id := range input.Permissions {
		permission, err := app.models.Permissions.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		permissions = append(permissions, *permission)
	}

	var role data.Role
	role.Name = input.Name
	role.Permissions = permissions

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

func (app *application) getAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	qs := r.URL.Query()
	p := &data.Paginate{
		Limit: app.readInt(qs, v, "limit", 10),
		Page:  app.readInt(qs, v, "page", 1),
	}
	if data.ValidatePaginate(v, p); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	roles, metadata, err := app.models.Roles.GetAll(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"roles": roles, "metadata": metadata}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) getRoleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"role": role}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateRolesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input roleDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var newPermissions []data.Permission
	for _, id := range input.Permissions {
		permission, err := app.models.Permissions.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
				return
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		newPermissions = append(newPermissions, *permission)
	}

	role.Name = input.Name
	role.Permissions = newPermissions

	if err := app.models.Roles.Update(role); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "resource updated"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteRolesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	role, err := app.models.Roles.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.models.Roles.Delete(role)
	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
