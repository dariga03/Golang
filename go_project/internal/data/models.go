package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support.
	Researchers interface {
		Insert(researcher *Researcher) error
		Get(id int64) (*Researcher, error)
		GetAll(name string, specialization string, filters Filters) ([]*Researcher, Metadata, error)
		Update(researcher *Researcher) error
		Delete(id int64) error
	}

	Expeditions interface {
		Insert(expedition *Expedition) error
		Get(id int64) (*Expedition, error)
		GetAll(title string, expeditionYear int, filters Filters) ([]*Expedition, Metadata, error)
		Update(expedition *Expedition) error
		Delete(id int64) error
		GetExpeditionsByResearcher(researcher_id int64, title string, expeditionYear int, filters Filters) ([]*Expedition, Metadata, error)
	}

	Artifacts interface {
		Insert(artifact *Artifact) error
		Get(id int64) (*Artifact, error)
		GetAll(title string, age int, filters Filters) ([]*Artifact, Metadata, error)
		Update(artifact *Artifact) error
		Delete(id int64) error
		GetArtifactsByResearcher(researcher_id int64, title string, age int, filters Filters) ([]*Artifact, Metadata, error)
	}
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

// Create a helper function which returns a Models instance containing the mock models
// only.
func NewModels(db *sql.DB) Models {
	return Models{
		Researchers: ResearcherModel{DB: db},
		Expeditions: ExpeditionModel{DB: db},
		Artifacts:   ArtifactModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
