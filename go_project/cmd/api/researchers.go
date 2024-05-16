package main

import (
	"errors"
	"fmt"
	"net/http"

	"goproject/internal/data"
	"goproject/internal/validator"
)

func (app *application) createResearcherHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string `json:"name"`
		Specialization string `json:"specialization"`
		Project        string `json:"project"`
	}

	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the input struct to a new Researcher struct.

	researcher := &data.Researcher{
		Name:           input.Name,
		Specialization: input.Specialization,
		Project:        input.Project,
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call the ValidateSong() function and return a response containing the errors if
	// any of the checks fail.
	if data.ValidateResearcher(v, researcher); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Call the Insert() method on our researchers model, passing in a pointer to the
	// validated song struct. This will create a record in the database and update the
	// song struct with the system-generated information.
	err = app.models.Researchers.Insert(researcher)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/researchers/%d", researcher.Id))
	// Write a JSON response with a 201 Created status code, the song data in the
	// response body, and the Location header.

	err = app.writeJSON(w, http.StatusCreated, envelope{"researcher": researcher}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// fmt.Fprintf(w, "%+v\n", input)
}

// Add a showMovieHandler for the "GET /v1/movies/:id" endpoint. For now, we retrieve
// the interpolated "id" parameter from the current URL and include it in a placeholder
// response.
func (app *application) showResearcherHandler(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, any interpolated URL parameters will be
	// stored in the request context. We can use the ParamsFromContext() function to
	// retrieve a slice containing these parameter names and values.
	//params := httprouter.ParamsFromContext(r.Context())

	// We can then use the ByName() method to get the value of the "id" parameter from
	// the slice. In our project all movies will have a unique positive integer ID, but
	// the value returned by ByName() is always a string. So we try to convert it to a
	// base 10 integer (with a bit size of 64). If the parameter couldn't be converted,
	// or is less than 1, we know the ID is invalid so we use the http.NotFound()
	// function to return a 404 Not Found response.
	id, err := app.readIDParam(r)
	if err != nil {
		// Use the new notFoundResponse() helper.
		app.notFoundResponse(w, r)
		return
	}

	// Call the Get() method to fetch the data for a specific movie. We also need to
	// use the errors.Is() function to check if it returns a data.ErrRecordNotFound
	// error, in which case we send a 404 Not Found response to the client.
	// movie, err := app.models.Movies.Get(id)
	researcher, err := app.models.Researchers.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Encode the struct to JSON and send it as the HTTP response.
	// err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	err = app.writeJSON(w, http.StatusOK, envelope{"researcher": researcher}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateResearcherHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the researcher ID from the URL.
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	// Fetch the existing movie record from the database, sending a 404 Not Found
	// response to the client if we couldn't find a matching record.
	// movie, err := app.models.Movies.Get(id)
	researcher, err := app.models.Researchers.Get(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Declare an input struct to hold the expected data from the client.

	var input struct {
		Name           string `json:"name"`
		Specialization string `json:"specialization"`
		Project        string `json:"project"`
	}

	// Read the JSON request body data into the input struct.
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from the request body to the appropriate fields of the movie
	// record.

	researcher.Name = input.Name
	researcher.Specialization = input.Specialization
	researcher.Project = input.Project

	// Validate the updated researcher record, sending the client a 422 Unprocessable Entity
	// response if any checks fail.
	v := validator.New()
	// if data.ValidateMovie(v, movie); !v.Valid() {
	// 	app.failedValidationResponse(w, r, v.Errors)
	// 	return
	// }
	if data.ValidateResearcher(v, researcher); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated researcher record to our new Update() method.
	err = app.models.Researchers.Update(researcher)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Write the updated researcher record in a JSON response.
	err = app.writeJSON(w, http.StatusOK, envelope{"researcher": researcher}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteResearcherHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the researcher ID from the URL.
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Delete the researcher from the database, sending a 404 Not Found response to the
	// client if there isn't a matching record.
	err = app.models.Researchers.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Return a 200 OK status code along with a success message.
	// err = app.writeJSON(w, http.StatusOK, envelope{"message": "researcher successfully deleted"}, nil)
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "researcher successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listResearchersHandler(w http.ResponseWriter, r *http.Request) {
	// Embed the new Filters struct.

	var input struct {
		Name           string `json:"name"`
		Specialization string `json:"specialization"`
		Project        string `json:"project"`
		data.Filters
	}

	// Initialize a new Validator instance.
	v := validator.New()

	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()

	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Name = app.readString(qs, "name", "")

	input.Specialization = app.readString(qs, "specialization", "")

	input.Project = app.readString(qs, "project", "")

	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "researcher_id")
	// Add the supported sort values for this endpoint to the sort safelist.
	// input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}
	input.Filters.SortSafelist = []string{"researcher_id", "name", "specialization", "project", "-researcher_id", "-name", "-specialization", "-project"}

	// Execute the validation checks on the Filters struct and send a response
	// containing the errors if necessary.
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	// Call the GetAll() method to retrieve the researchers, passing in the various filter
	// parameters.
	// Accept the metadata struct as a return value.
	// researchers, metadata, err := app.models.Researchers.GetAll(input.Name, input.Specialization, input.Project)
	researchers, metadata, err := app.models.Researchers.GetAll(input.Name, input.Specialization, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the researcher data.
	// Include the metadata in the response envelope.
	// err = app.writeJSON(w, http.StatusOK, envelope{"researchers": researchers, "metadata": metadata}, nil)
	err = app.writeJSON(w, http.StatusOK, envelope{"researchers": researchers, "metadata": metadata}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
