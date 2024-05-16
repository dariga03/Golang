package data

import (
	"goproject/internal/validator"
	"context"
	"database/sql"
	"errors"
	"time"
	"fmt"

	//"github.com/lib/pq"
)


type Researcher struct{
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	Project        string `json:"project"`
} 

func ValidateResearcher(v *validator.Validator, researcher *Researcher){
	v.Check(researcher.Name != "", "name", "must be provided")
	v.Check(researcher.Project != "", "project", "must be greater than 0")
}
/*
func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
*/

// Define a MovieModel struct type which wraps a sql.DB connection pool.

type ResearcherModel struct {
	DB *sql.DB
}


// Add a placeholder method for inserting a new record in the researchers table.
func (s ResearcherModel) Insert(researcher *Researcher) error {
	// Define the SQL query for inserting a new record in the researchers table and returning
	// the system-generated data.
	query := `
		INSERT INTO researcher(name, specialization, project)
		VALUES ($1, $2, $3)
		RETURNING researcher_id, name, specialization, project;`

	// Create an args slice containing the values for the placeholder parameters from
	// the reseracher struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{researcher.Name, researcher.Specialization, researcher.Project}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the researcher struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&researcher.Id, &researcher.Name, &researcher.Specialization, &researcher.Project)
}

// Add a placeholder method for fetching a specific record from the researchers table.
func (s ResearcherModel) Get(id int64) (*Researcher, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT researcher_id, name, specialization, project
		FROM researcher
		WHERE researcher_id = $1;`

	var researcher Researcher
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&researcher.Id, &researcher.Name, &researcher.Specialization, &researcher.Project)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &researcher, nil
}

// Add a placeholder method for updating a specific record in the researchers table.
func (s ResearcherModel) Update(researcher *Researcher) error {
	query := `
		UPDATE researcher
		SET name = $1, specialization = $2, project = $3
		WHERE researcher_id = $4
		RETURNING researcher_id, name, specialization, project;
		`

	args := []interface{}{researcher.Name, researcher.Specialization, researcher.Project, researcher.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&researcher.Id, &researcher.Name, &researcher.Specialization, &researcher.Project)
}


func (s ResearcherModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the researcher ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM researcher
		WHERE researcher_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	// If no rows were affected, we know that the researchers table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}


// Create a new GetAll() method which returns a slice of researchers. Although we're not
// using them right now, we've set this up to accept the various filter parameters as
// arguments.
func (s ResearcherModel) GetAll(name string, specialization string, filters Filters) ([]*Researcher, Metadata, error) {
	// Construct the SQL query to retrieve all researcher records.
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), researcher_id, name, specialization, project
		FROM researcher
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', specialization) @@ plainto_tsquery('simple', $2) OR $2 = '')
		ORDER BY %s %s, researcher_id
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())
	

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{name, specialization, filters.limit(), filters.offset()}

	// Use QueryContext() to execute the query. This returns a sql.Rows resultset
	// containing the result.
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	// Importantly, defer a call to rows.Close() to ensure that the resultset is closed
	// before GetAll() returns.
	defer rows.Close()
	// Declare a totalRecords variable.
	totalRecords := 0
	// Initialize an empty slice to hold the researcher data.
	researchers := []*Researcher{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var researcher Researcher
		// Scan the values from the row into the Researcher struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&researcher.Id,
			&researcher.Name,
			&researcher.Specialization,
			&researcher.Project,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Researcher struct to the slice.
		researchers = append(researchers, &researcher)
	}
	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	// Generate a Metadata struct, passing in the total record count and pagination
	// parameters from the client.
	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	// If everything went OK, then return the slice of researchers.
	return researchers, metadata, nil
}
