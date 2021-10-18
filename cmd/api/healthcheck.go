package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	out := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	if err := app.writeJSON(w, http.StatusOK, out, nil); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
