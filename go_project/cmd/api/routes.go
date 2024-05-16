package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Convert the notFoundResponse() helper to a http.Handler using the
	// http.HandlerFunc() adapter, and then set it as the custom error handler for 404
	// Not Found responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	// Likewise, convert the methodNotAllowedResponse() helper to a http.Handler and set
	// it as the custom error handler for 405 Method Not Allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the relevant methods, URL patterns and handler functions for our
	// endpoints using the HandlerFunc() method. Note that http.MethodGet and
	// http.MethodPost are constants which equate to the strings "GET" and "POST"
	// respectively.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/researchers", app.requirePermission("read", app.listResearchersHandler))
	router.HandlerFunc(http.MethodPost, "/v1/researchers", app.requirePermission("write", app.createResearcherHandler))
	router.HandlerFunc(http.MethodGet, "/v1/researchers/:id", app.requirePermission("read", app.showResearcherHandler))
	router.HandlerFunc(http.MethodPut, "/v1/researchers/:id", app.requirePermission("write", app.updateResearcherHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/researchers/:id", app.requirePermission("write", app.deleteResearcherHandler))

	router.HandlerFunc(http.MethodGet, "/v1/expeditions", app.requirePermission("read", app.listExpeditionsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/expeditions", app.requirePermission("write", app.createExpeditionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/expeditions/:id", app.requirePermission("read", app.showExpeditionHandler))
	router.HandlerFunc(http.MethodPut, "/v1/expeditions/:id", app.requirePermission("write", app.updateExpeditionHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/expeditions/:id", app.requirePermission("write", app.deleteExpeditionHandler))
	router.HandlerFunc(http.MethodGet, "/v1/researchers/:id/expeditions", app.requirePermission("read", app.getExpeditionsByResearcherHandler))

	router.HandlerFunc(http.MethodGet, "/v1/artifacts", app.requirePermission("read", app.listArtifactsHandler))
	router.HandlerFunc(http.MethodPost, "/v1/artifacts", app.requirePermission("write", app.createArtifactHandler))
	router.HandlerFunc(http.MethodGet, "/v1/artifacts/:id", app.requirePermission("read", app.showArtifactHandler))
	router.HandlerFunc(http.MethodPut, "/v1/artifacts/:id", app.requirePermission("write", app.updateArtifactHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/artifacts/:id", app.requirePermission("write", app.deleteArtifactHandler))
	router.HandlerFunc(http.MethodGet, "/v1/researchers/:id/artifacts", app.requirePermission("read", app.getArtifactsByResearcherHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/login", app.createAuthenticationTokenHandler)

	// Return the httprouter instance.
	// return router
	return app.authenticate(router)

}
