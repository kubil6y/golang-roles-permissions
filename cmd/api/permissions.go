package main

import (
	"errors"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

func (app *application) createPermissionHandler(w http.ResponseWriter, r *http.Request) {
	var input permissionDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	var permission data.Permission
	input.populate(&permission)

	if err := app.models.Permissions.Insert(&permission); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("name", "a permission with that name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"permission": permission}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllPermissionsHandler(w http.ResponseWriter, r *http.Request) {
	permissions, err := app.models.Permissions.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"permissions": permissions}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getPermissionHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

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

	e := envelope{"permission": permission}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updatePermissionsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input permissionDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

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

	input.populate(permission)

	if err := app.models.Permissions.Update(permission); err != nil {
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

func (app *application) grantPermissionsToRolesHandler(w http.ResponseWriter, r *http.Request) {
	var input grantPermissionsToRolesDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	role, err := app.models.Roles.GetByID(input.RoleID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// TODO
	// 1- logic messed up first check in rolespermissions then look for permissions

	// client can send [1,1,1,2] for permission ids, extra DB calls
	setOfIncomingPermissionIDs := app.intSliceToSet(input.Permissions)

	var newPermissions []*data.Permission
	var notFoundPermissionIDs []int64
	var alreadyExistingPermissionIDs []int64

	// check if permission already exists in role's permissions
	for _, existing := range role.Permissions {
		for _, incoming := range newPermissions {
			if existing.ID == incoming.ID {
				alreadyExistingPermissionIDs = append(alreadyExistingPermissionIDs, existing.ID)
			}
		}
	}

	if len(alreadyExistingPermissionIDs) > 0 {
		// take these out from setOfIncomingPermissionIDs
		// TODO leftoff
	}

	// checking permissions exists or not
	for _, id := range setOfIncomingPermissionIDs {
		permission, err := app.models.Permissions.GetByID(id)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				notFoundPermissionIDs = append(notFoundPermissionIDs, id)
			} else {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		if permission != nil {
			newPermissions = append(newPermissions, permission)
		}
	}

	// adding new permissions
	for _, np := range newPermissions {
		role.Permissions = append(role.Permissions, *np)
	}

	// saving role
	app.models.Roles.DB.Save(role)

	e := envelope{"role": role}

	if len(notFoundPermissionIDs) > 0 {
		e["not_found_permission_ids"] = notFoundPermissionIDs
	}

	if len(alreadyExistingPermissionIDs) > 0 {
		e["already_existing_permission_ids"] = alreadyExistingPermissionIDs
	}

	var out map[string]interface{}
	if len(notFoundPermissionIDs) > 0 || len(alreadyExistingPermissionIDs) > 0 {
		out = app.outERR(e)
	} else {
		out = app.outOK(e)
	}

	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
