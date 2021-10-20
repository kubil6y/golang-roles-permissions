package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.isAdmin(app.healthCheckHandler))

	router.HandlerFunc(http.MethodPost, "/v1/admin/permissions", app.isAdmin(app.createPermissionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/permissions", app.isAdmin(app.getAllPermissionsHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/permissions/:id", app.isAdmin(app.getPermissionHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/permissions/:id", app.isAdmin(app.updatePermissionsHandler))

	router.HandlerFunc(http.MethodPost, "/v1/admin/roles", app.isAdmin(app.createRoleHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles", app.isAdmin(app.getAllRolesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/admin/roles/:id", app.isAdmin(app.getRoleHandler))
	router.HandlerFunc(http.MethodPut, "/v1/admin/roles/:id", app.isAdmin(app.updateRolesHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

// NOTE when trying to access an invalid or expired token,
// clients wont be able to login. remove it on client side.
