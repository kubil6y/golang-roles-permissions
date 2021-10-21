package main

/*
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

	e := envelope{}
	var newPermissions []data.Permission
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
		e["already_existing_permission_ids"] = alreadyExistingPermissionIDs
		out := app.outERR(e)
		if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// checking permissions exists or not
	for _, id := range input.Permissions {
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
			newPermissions = append(newPermissions, *permission)
		}
	}

	if len(notFoundPermissionIDs) > 0 {
		e["not_found_permission_ids"] = notFoundPermissionIDs
		out := app.outERR(e)
		if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// adding new permissions
	role.Permissions = append(role.Permissions, newPermissions...)
	// saving new role.Permissions
	app.models.Roles.DB.Save(role)

	e = envelope{"role": role}
	out := app.outOK(e)
	if err := app.writeJSON(w, http.StatusCreated, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
*/
