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

// UpdateVersionLatestCommit updates the latest commit of the specified version
func (r PgVersionRepository) UpdateVersionLatestCommit(version int64, commit string) error {

	// update field
	stmt, err := r.Db.Prepare("UPDATE version SET latestCommit=$2 where id=$1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(version, commit)
	if err != nil {
		return err
	}

	return nil
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
    			status VARCHAR(16) NOT NULL DEFAULT 'draft',
    			latestCommit CHAR(40),
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
	stmt, err := r.Db.Prepare("SELECT articleid, id, title, status, latestCommit FROM version WHERE id=$1")

	if err != nil {
		return entities.Version{}, err
	}
	row := stmt.QueryRow(version)

	// Extract results
	var entity entities.Version
	err = row.Scan(&entity.ArticleID, &entity.Id, &entity.Title, &entity.Status, &entity.LatestCommitID)
	if err != nil {
		return entities.Version{}, err
	}

	return entity, nil
}

// GetVersionsByArticle gets the version entities related to a specific article, links the owners
func (r PgVersionRepository) GetVersionsByArticle(article int64) ([]entities.Version, error) {

	// Prepare and execute query
	stmt, err := r.Db.Prepare(`SELECT id, title, status, versionowners.email, latestCommit
		FROM version INNER JOIN versionowners ON version.id = versionowners.versionid
		WHERE articleid=$1`)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(article)
	if err != nil {
		return nil, err
	}

	// Extract results into map,
	// Necessary because the same version might appear in multiple rows if it has multiple owners
	versions := make(map[int64]entities.Version)
	for rows.Next() {
		// Read the current row
		row := entities.Version{}
		var email string
		if err := rows.Scan(&row.Id, &row.Title, &row.Status, &email, &row.LatestCommitID); err != nil {
			return nil, err
		}

		// Insert new version into map, or append email to existing
		if version, ok := versions[row.Id]; ok {
			// Exists
			version.Owners = append(version.Owners, email)
			versions[row.Id] = version
		} else {
			// Add new, with list of just one email
			row.Owners = []string{email}
			versions[row.Id] = row
		}
	}

	// Turn map into go slice
	var result []entities.Version
	for _, value := range versions {
		result = append(result, value)
	}

	return result, nil
}

// GetOwners gets a string of owner emails, belonging to the specified version
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

// CheckIfOwner returns directly with a query whether the specified user owns an article version
func (r PgVersionRepository) CheckIfOwner(version int64, email string) (bool, error) {
	// Prepare and execute query
	stmt, err := r.Db.Prepare("SELECT 1 FROM versionOwners WHERE versionid=$1 AND email=$2")
	if err != nil {
		return false, err
	}

	rows, err := stmt.Query(version, email)
	if rows.Next() {
		// there is (at least) one row, so the specified email is an owner of the version
		return true, nil
	} else {
		// no rows match, so the specified email is not an owner of the version
		return false, nil
	}
}
