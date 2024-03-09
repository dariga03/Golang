package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dariga03/go-app/pkg/abr-plus/model"
	"github.com/gorilla/mux"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Creating a new artifact
func (app *application) createArtifactHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Id             string `json:"id"`
		Title          string `json:"title"`
		Description    string `json:"description"`
		Age            int    `json:"age"`
		Location       string `json:"location"`
		Researchers_id int    `json:"researchers_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	artifact := &model.Artifact{
		Id:             input.Id,
		Title:          input.Title,
		Description:    input.Description,
		Age:            input.Age,
		Location:       input.Location,
		Researchers_id: input.Researchers_id,
	}

	err = app.models.Artifacts.Insert(artifact)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, artifact)
}

// Getting a specific song by id
func (app *application) getArtifactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["artifactId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	artifact, err := app.models.Artifacts.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	app.respondWithJSON(w, http.StatusOK, artifact)
}

// Updating a artifact by id
func (app *application) updateArtifactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["artifactId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	artifact, err := app.models.Artifacts.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "404 Not Found")
		return
	}

	var input struct {
		Id             *string `json:"id"`
		Title          *string `json:"title"`
		Description    *string `json:"description"`
		Age            *int    `json:"age"`
		Location       *string `json:"location"`
		Researchers_id *int    `json:"researchers_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		artifact.Title = *input.Title
	}

	if input.Description != nil {
		artifact.Description = *input.Description
	}

	if input.Age != nil {
		artifact.Age = *input.Age
	}

	if input.Location != nil {
		artifact.Location = *input.Location
	}

	if input.Researchers_id != nil {
		artifact.Researchers_id = *input.Researchers_id
	}

	err = app.models.Artifacts.Update(artifact)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, artifact)
}

// Deleting a artifact by id
func (app *application) deleteArtifactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["artifactId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid artifact ID")
		return
	}

	err = app.models.Artifacts.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (app *application) readJSON(_ http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
