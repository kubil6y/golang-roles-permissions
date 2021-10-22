package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
	"github.com/kubil6y/myshop-go/internal/validator"
)

// TODO give user default permissions after creating roles and shit
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
	user.SetPassword(input.Password)
	user.IsActivated = false

	if err := app.models.Users.Insert(&user); err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateRecord):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"user": user}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	v := validator.New()
	p := &data.Paginate{
		Limit: app.readInt(qs, v, "limit", 10),
		Page:  app.readInt(qs, v, "page", 1),
	}

	if data.ValidatePaginate(v, p); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	users, metadata, err := app.models.Users.GetAll(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{
		"users":    users,
		"metadata": metadata,
	}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	e := envelope{"user": user}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input updateUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := input.populate(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.models.Users.Update(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.models.Users.Delete(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateUserOwnHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input updateUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := input.populate(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.models.Users.Update(user); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	me := app.contextGetUser(r)
	e := envelope{"user": me}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) grantRoleToUserHandler(w http.ResponseWriter, r *http.Request) {
	var input roleToUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	targetUser, err := app.models.Users.GetByIDWithRolesAndPermissions(input.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var inputRoles []data.Role // roles in dto
	for _, roleID := range input.RoleIDs {
		role, err := app.models.Roles.GetByID(roleID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		inputRoles = append(inputRoles, *role)
	}

	targetUser.Roles = append(targetUser.Roles, inputRoles...)

	if err := app.models.Users.Update(targetUser); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) revokeRoleToUserHandler(w http.ResponseWriter, r *http.Request) {
	var input roleToUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	targetUser, err := app.models.Users.GetByIDWithRolesAndPermissions(input.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var inputRoles []data.Role // roles in dto
	for _, roleID := range input.RoleIDs {
		role, err := app.models.Roles.GetByID(roleID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		inputRoles = append(inputRoles, *role)
	}

	doesRoleExist := func(list []data.Role, role data.Role) bool {
		for _, v := range list {
			if v.ID == role.ID {
				return true
			}
		}
		return false
	}

	var newRoles []data.Role

	for _, existingRole := range targetUser.Roles {
		if !doesRoleExist(inputRoles, existingRole) {
			newRoles = append(newRoles, existingRole)
		}
	}

	targetUser.Roles = newRoles

	if err := app.models.Users.UpdateRoles(targetUser); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) grantPermissionToUserHandler(w http.ResponseWriter, r *http.Request) {
	var input permissionToUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByIDWithRolesAndPermissions(input.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var inputPermissions []data.Permission
	for _, permissionID := range input.PermissionIDs {
		permission, err := app.models.Permissions.GetByID(permissionID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		inputPermissions = append(inputPermissions, *permission)
	}

	for _, ip := range inputPermissions {
		if ContainsPermission(user.GrantedPermissions, ip) {
			continue
		}
		user.GrantedPermissions = append(user.GrantedPermissions, ip)
	}

	for _, v := range user.GrantedPermissions {
		fmt.Println("granted id", v.ID)
	}

	var filteredRevokedPermissions []data.Permission
	for _, rp := range user.RevokedPermissions {
		if ContainsPermission(inputPermissions, rp) {
			continue
		}
		filteredRevokedPermissions = append(filteredRevokedPermissions, rp)
	}
	user.RevokedPermissions = filteredRevokedPermissions

	err = app.models.Users.UpdateGrantedPermissions(user)
	err = app.models.Users.UpdateRevokedPermissions(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) revokePermissionToUserHandler(w http.ResponseWriter, r *http.Request) {
	var input permissionToUserDTO
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if input.validate(v); !v.IsValid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByIDWithRolesAndPermissions(input.UserID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var inputPermissions []data.Permission
	for _, permissionID := range input.PermissionIDs {
		permission, err := app.models.Permissions.GetByID(permissionID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		inputPermissions = append(inputPermissions, *permission)
	}

	// deciding new granted permissions TODO
	var filteredGrantedPermissions []data.Permission
	for _, grantedPerm := range user.GrantedPermissions {
		if !ContainsPermission(inputPermissions, grantedPerm) {
			filteredGrantedPermissions = append(filteredGrantedPermissions, grantedPerm)
		}
	}
	user.GrantedPermissions = filteredGrantedPermissions

	// deciding new revoked permissions TODO
	for _, ip := range inputPermissions {
		if ContainsPermission(user.RevokedPermissions, ip) {
			continue
		}
		user.RevokedPermissions = append(user.RevokedPermissions, ip)
	}

	err = app.models.Users.UpdateGrantedPermissions(user)
	err = app.models.Users.UpdateRevokedPermissions(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	e := envelope{"message": "success"}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusAccepted, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getUserRolesAndPermissions(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByIDWithRolesAndPermissions(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	e := envelope{"user": user}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
