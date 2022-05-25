package postgres

import (
	"database/sql"
	"mainServer/entities"
)

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

func (r PgRequestRepository) CreateRequest(req entities.Request) (entities.Request, error) {

	// store request entity
	stmt, err := r.Db.Prepare("INSERT INTO request (id, articleID, sourceVersionID, sourceHistoryID, targetVersionID, targetHistoryID) VALUES (DEFAULT,$1,$2,$3,$4,$5) RETURNING id")
	if err != nil {
		return entities.Request{}, err
	}
	row := stmt.QueryRow(req.ArticleID, req.SourceVersionID, req.SourceHistoryID, req.TargetVersionID, req.TargetHistoryID)

	// Because QueryRow instead of Exec was used, with the RETURNING keyword,
	// The generated id can be retrieved
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return entities.Request{}, err
	}

	req.RequestID = id
	return req, nil
}

// creates request table, storing history/commit IDs as fixed-length,
// as they are always sha-1 hashes of 40 hex digits long
func (r PgRequestRepository) createRequestTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS request (
    	id SERIAL PRIMARY KEY,
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
