package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
)

type PgVersionRepository struct {
	Db *sql.DB
}

func (r PgVersionRepository) GetVersion(version int64) (entities.Version, error) {
	// get version without owners
	entity, err := r.getVersion(version)
	if err != nil {
		return entities.Version{}, err
	}

	// add owners
	owners, err := r.getOwners(version)
	if err != nil {
		return entities.Version{}, err
	}

	entity.Owners = owners
	return entity, nil

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

	// add entry in version table
	id, err := r.createVersion(version.ArticleID, version.Title)
	if err != nil {
		return entities.Version{}, err
	}

	// link owners in versionOwner table
	err = r.linkOwners(id, version.Owners)
	if err != nil {
		return entities.Version{}, err
	}

	version.Id = id
	return version, nil
}

// createVersion adds a single row to the version table, returns generated id.
// Does not link owners.
func (r PgVersionRepository) createVersion(article int64, title string) (int64, error) {

	// store version entity
	stmt, err := r.Db.Prepare("INSERT INTO version (id, articleID, title) VALUES (DEFAULT,$1,$2) RETURNING id")
	if err != nil {
		return -1, err
	}
	row := stmt.QueryRow(article, title)

	// Because QueryRow instead of Exec was used, with the RETURNING keyword,
	// The generated id can be retrieved
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

// linkOwners inserts rows to specify the owners of a version.
// note: does not delete rows, so will not remove existing owners, if they are excluded in the function parameter.
func (r PgVersionRepository) linkOwners(version int64, owners []string) error {
	query := `INSERT INTO versionOwners (versionID, email) VALUES `

	// add a parametrized value list to the query dynamically, using the length of owners
	// (<versionID>, $1), (<versionID>, $2), ...
	// also converts owners from []string to []any (required for stmt.Exec)
	// implementation inspired by https://stackoverflow.com/a/51132288
	var values []any
	for i, s := range owners {
		values = append(values, s)
		query = query + fmt.Sprintf("(%d,$%d),", version, i+1)
	}
	query = query[:len(query)-1] // remove trailing comma

	// prepare query
	stmt, err := r.Db.Prepare(query)
	if err != nil {
		return err
	}

	// execute the statement by inserting the owner emails in the $i's
	_, err = stmt.Exec(values...)
	return err
}

func (r PgVersionRepository) createVersionTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS version (
    			id SERIAL PRIMARY KEY,
    			articleID int NOT NULL,
    			title VARCHAR(255) NOT NULL,
    			FOREIGN KEY (articleID) REFERENCES article(id)    			
    )`)
	return err
}

func (r PgVersionRepository) createVersionOwnersTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS versionOwners (
    			email VARCHAR(255) NOT NULL,
    			versionID int NOT NULL,
    			FOREIGN KEY (email) REFERENCES users(email),
    			FOREIGN KEY (versionID) REFERENCES version(id)    			 
    )`)
	return err
}

// getVersion gets the version entity from the database, but doesn't link the owners
func (r PgVersionRepository) getVersion(version int64) (entities.Version, error) {
	// Prepare and execute query
	stmt, err := r.Db.Prepare("SELECT articleid, id, title FROM version WHERE id=$1")
	if err != nil {
		return entities.Version{}, err
	}
	row := stmt.QueryRow(version)

	// Extract results
	var entity entities.Version
	err = row.Scan(&entity.ArticleID, &entity.Id, &entity.Title)
	if err != nil {
		return entities.Version{}, err
	}

	return entity, nil
}

// getOwners gets a string of owner emails, belonging to the specified version
func (r PgVersionRepository) getOwners(version int64) ([]string, error) {
	// Prepare and execute query
	stmt, err := r.Db.Prepare("SELECT email FROM versionOwners WHERE versionid=$1")
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(version)
	if err != nil {
		return nil, err
	}

	// Extract results
	var owners []string
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		owners = append(owners, email)
	}

	return owners, nil
}
