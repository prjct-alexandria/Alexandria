package postgres

import "database/sql"

type PgRequestRepository struct {
	Db *sql.DB
}

// NewPgRequestRepository creates a new repository and tables in the database, if necessary.
func NewPgRequestRepository(db *sql.DB) PgRequestRepository {
	repo := PgRequestRepository{Db: db}
	err := repo.createRequestTable()
	if err != nil {
		panic(err)
	}
	return repo
}

// creates request table, storing history/commit IDs as fixed-length,
// as they are always sha-1 hashes of 40 hex digits long
func (r PgRequestRepository) createRequestTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS request (
    	requestID SERIAL PRIMARY KEY,
    	articleID INT NOT NULL,
    	sourceVersionID INT NOT NULL,
    	sourceHistoryID NCHAR(40) NOT NULL,
    	targetVersionID INT NOT NULL,
    	targetHistoryID NCHAR(40) NOT NULL,
    	state VARCHAR(16) NOT NULL DEFAULT 'pending',
        FOREIGN KEY (articleID) REFERENCES article(id),
        FOREIGN KEY (sourceVersionID) REFERENCES version(id),
        FOREIGN KEY (targetVersionID) REFERENCES version(id)
    )`)
	return err
}
