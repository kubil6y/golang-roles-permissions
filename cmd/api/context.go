package main

import (
	"context"
	"net/http"

	"github.com/kubil6y/myshop-go/internal/data"
)

// contextKey is a type for avoiding name clashes.
type contextKey string

const userContextKey = contextKey("user")

func (app *application) setUserContext(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	// by default Value() returns interface{}, so we use type assertion.
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
