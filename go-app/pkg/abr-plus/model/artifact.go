package model

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Artifact struct {
	Id             string `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Age            int    `json:"age"`
	Location       string `json:"location"`
	Researchers_id int    `json:"researchers_id"`
}

type ArtifactModel struct {
	DB       *sql.DB
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

// Insert a new artifact item into the database.
func (m ArtifactModel) Insert(artifact *Artifact) error {
	query := `
		INSERT INTO artifact (artifact_id, title, description, age, location, researshers_id) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING artifact_id, title, description, age, location, researshers_id;
		`
	args := []interface{}{artifact.Id, artifact.Title, artifact.Description, artifact.Age, artifact.Location, artifact.Researchers_id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&artifact.Id, &artifact.Title, &artifact.Description, &artifact.Age, &artifact.Location, &artifact.Researchers_id)
}

// Getting a artifact by id
func (s ArtifactModel) Get(id int) (*Artifact, error) {
	query := `
		SELECT artifact_id, title, description, age, location, researchers_id
		FROM artifact
		WHERE artifact_id = $1;
		`
	var artifact Artifact
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&artifact.Id, &artifact.Title, &artifact.Description, &artifact.Age, &artifact.Location, &artifact.Researchers_id)
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}

// Updating a artifact by id
func (s ArtifactModel) Update(artifact *Artifact) error {
	query := `
		UPDATE artifact
		SET title = $1, age = $2
		WHERE artifact_id = $3
		RETURNING artifact_id, title, description, age, location, researchers_id;
		`
	args := []interface{}{artifact.Title,  artifact.Age, artifact.Id}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return s.DB.QueryRowContext(ctx, query, args...).Scan(&artifact.Id, &artifact.Title, &artifact.Description, &artifact.Age, &artifact.Location, &artifact.Researchers_id)
}

// Deleting a artifact by id
func (s ArtifactModel) Delete(id int) error {
	query := `
		DELETE FROM artifact
		WHERE artifact_id = $1
		`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := s.DB.ExecContext(ctx, query, id)
	return err
}
