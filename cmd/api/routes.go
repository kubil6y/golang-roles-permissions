package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requirePermission("perm100", (app.healthCheckHandler)))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users", app.getAllUsersHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/me", app.updateUserOwnHandler)
	router.HandlerFunc(http.MethodGet, "/v1/profile", app.getProfileHandler)

	router.HandlerFunc(http.MethodPost, "/v1/admin/permissions", app.requirePermission("admin", app.createPermissionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/permissions", app.requirePermission("admin", app.getAllPermissionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/permissions/:id", app.requirePermission("admin", app.getPermissionHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/permissions/:id", app.requirePermission("admin", app.updatePermissionsHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/permissions/:id", app.requirePermission("admin", app.deletePermissionsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.requirePermission("admin", app.createRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.requirePermission("admin", app.getAllRolesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.requirePermission("admin", app.getRoleHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.requirePermission("admin", app.updateRolesHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/roles/:id", app.requirePermission("admin", app.deleteRolesHandler))

	router.HandlerFunc(http.MethodGet, "/v1/admin/users/access/:id", app.requirePermission("admin", app.getUserRolesAndPermissions))
	router.HandlerFunc(http.MethodPost, "/v1/admin/users/grant-role", app.requirePermission("admin", app.grantRoleToUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/admin/users/revoke-role", app.requirePermission("admin", app.revokeRoleToUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/admin/users/grant-permission", app.requirePermission("admin", app.grantPermissionToUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/admin/users/revoke-permission", app.requirePermission("admin", app.revokePermissionToUserHandler))

	router.HandlerFunc(http.MethodPatch, "/v1/admin/users/:id", app.requirePermission("admin", app.updateUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/admin/users/:id", app.requirePermission("admin", app.deleteUserHandler))

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

// NOTE when trying to access an invalid or expired token,
// clients wont be able to login. remove it on client side.
