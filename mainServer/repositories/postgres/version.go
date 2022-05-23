package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"strings"
)

type PgVersionRepository struct {
	Db *sql.DB
}

func NewPgVersionRepository(db *sql.DB) PgVersionRepository {
	repo := PgVersionRepository{db}

	err := repo.createVersionTable()
	if err != nil {
		panic(err)
	}

	err = repo.createVersionOwnersTable()
	if err != nil {
		panic(err)
	}
	return repo
}

func (r PgVersionRepository) CreateVersion(version entities.Version) (entities.Version, error) {
	title := version.Title
	owners := version.Owners

	// store version entity
	result, err := r.Db.Exec(fmt.Sprintf("INSERT INTO version (id, title) VALUES (DEFAULT, %s)", version.Title))
	if err != nil {
		return entities.Version{}, err
	}

	// read the generated id to return, should always be here if the query itself returned no error
	id, _ := result.LastInsertId()

	// create a string with pairs like "(versionID, email1), (versionID, email2), (versionID, email3)"
	center := strings.Join(owners, fmt.Sprintf("),(%d,", id))
	pairs := fmt.Sprintf("(%s, %s)", id, center)

	// store owners of this version, by inserting the created pairs
	_, err = r.Db.Exec(fmt.Sprintf("INSERT INTO versionOwners (id, title) VALUES %s", pairs))
	if err != nil {
		return entities.Version{}, err
	}

	return entities.Version{Id: id, Title: title, Owners: owners}, nil
}

func (r PgVersionRepository) createVersionTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS version (
    			id SERIAL PRIMARY KEY,
    			title VARCHAR(255) NOT NULL    			
    )`)
	return err
}

func (r PgVersionRepository) createVersionOwnersTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS versionOwners (
    			email VARCHAR(255) NOT NULL,
    			versionID VARCHAR(255) NOT NULL,
    			FOREIGN KEY (email) REFERENCES users(email),
    			FOREIGN KEY (versionID) REFERENCES version(id)    			 
    )`)
	return err
}
