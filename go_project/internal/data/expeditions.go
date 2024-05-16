package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goproject/internal/validator"
	"time"
	//"github.com/lib/pq"
)

type Expedition struct {
	Id             int    `json:"id"`
	Title          string `json:"title"`
	ExpeditionYear int    `json:"expeditionYear"`
	Researcher_id  int    `json:"researcher_id"`
}

func ValidateExpedition(v *validator.Validator, expedition *Expedition) {
	v.Check(expedition.Title != "", "title", "must be provided")	
	v.Check(expedition.ExpeditionYear  > 0, "expeditionYear", "must be greater than 0")
	v.Check(expedition.Researcher_id > 0, "researcher_id", "must be greater than 0")
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

type ExpeditionModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the researchers table.
func (s ExpeditionModel) Insert(expedition *Expedition) error {
	// Define the SQL query for inserting a new record in the researchers table and returning
	// the system-generated data.
	query := `
		INSERT INTO expedition( title, expeditionYear, researcher_id)
		VALUES ($1, $2, $3)
		RETURNING expedition_id, title, expeditionYear, researcher_id;`

	// Create an args slice containing the values for the placeholder parameters from
	// the reseracher struct. Declaring this slice immediately next to our SQL query helps to
	// make it nice and clear *what values are being used where* in the query.
	args := []interface{}{expedition.Title, expedition.ExpeditionYear, expedition.Researcher_id}

	// Use the QueryRow() method to execute the SQL query on our connection pool,
	// passing in the args slice as a variadic parameter and scanning the system-
	// generated id, created_at and version values into the researcher struct.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&expedition.Id, &expedition.Title, &expedition.ExpeditionYear, &expedition.Researcher_id)
}

// Add a placeholder method for fetching a specific record from the researchers table.
func (s ExpeditionModel) Get(id int64) (*Expedition, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	// Retrieve a specific menu item based on its ID.
	query := `
		SELECT expedition_id, title, expeditionYear, researcher_id
		FROM expedition
		WHERE expedition_id = $1;`

	var expedition Expedition
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&expedition.Id, &expedition.Title, &expedition.ExpeditionYear, &expedition.Researcher_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &expedition, nil
}

// Add a placeholder method for updating a specific record in the researchers table.
func (s ExpeditionModel) Update(expedition *Expedition) error {
	query := `
		UPDATE expedition
		SET title = $1, expeditionYear = $2, researcher_id = $3
		WHERE expedition_id = $4
		RETURNING expedition_id, title, expeditionYear, researcher_id;
		`

	args := []interface{}{expedition.Title, expedition.ExpeditionYear, expedition.Researcher_id, expedition.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&expedition.Id, &expedition.Title, &expedition.ExpeditionYear, &expedition.Researcher_id)
}

func (s ExpeditionModel) Delete(id int64) error {
	// Return an ErrRecordNotFound error if the researcher ID is less than 1.
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM expedition
		WHERE expedition_id = $1
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
func (s ExpeditionModel) GetAll(title string, expeditionYear int, filters Filters) ([]*Expedition, Metadata, error) {
	// Construct the SQL query to retrieve all researcher records.
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), expedition_id, title, expeditionYear, researcher_id
		FROM expedition
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (expeditionYear = $2 OR $2 = 1)
		ORDER BY %s %s, expedition_id
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	// Create a context with a 3-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []interface{}{title, expeditionYear, filters.limit(), filters.offset()}

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
	expeditions := []*Expedition{}

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		// Initialize an empty Movie struct to hold the data for an individual movie.
		var expedition Expedition
		// Scan the values from the row into the Researcher struct. Again, note that we're
		// using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords, // Scan the count from the window function into totalRecords.
			&expedition.Id,
			&expedition.Title,
			&expedition.ExpeditionYear,
			&expedition.Researcher_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		// Add the Researcher struct to the slice.
		expeditions = append(expeditions, &expedition)
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
	return expeditions, metadata, nil
}



func (s ExpeditionModel) GetExpeditionsByResearcher(id int64, title string, expeditionYear int, filters Filters) ([]*Expedition, Metadata, error) {
	query :=  fmt.Sprintf(`
		SELECT count(*) OVER(), expedition_id, title, expeditionYear, researcher_id
		FROM expedition
		WHERE (researcher_id = $1)
		AND (to_tsvector('simple', title) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (expeditionYear = $3 OR $3 = 1)
		ORDER BY %s %s, expedition_id
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()


	args := []interface{}{id, title, expeditionYear, filters.limit(), filters.offset()}

	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	
	defer rows.Close()
	
	totalRecords := 0
	
	expeditions := []*Expedition{}

	
	for rows.Next() {
		
		var expedition Expedition
		
		err := rows.Scan(
			&totalRecords, 
			&expedition.Id,
			&expedition.Title,
			&expedition.ExpeditionYear,
			&expedition.Researcher_id,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		
		expeditions = append(expeditions, &expedition)
	}
	
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}


	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return expeditions, metadata, nil
}
