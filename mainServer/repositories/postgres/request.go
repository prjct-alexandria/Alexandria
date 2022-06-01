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
	stmt, err := r.Db.Prepare("INSERT INTO request (articleID, sourceVersionID, sourceHistoryID, targetVersionID, targetHistoryID) VALUES ($1,$2,$3,$4,$5) RETURNING id,state")
	if err != nil {
		return entities.Request{}, err
	}
	row := stmt.QueryRow(req.ArticleID, req.SourceVersionID, req.SourceHistoryID, req.TargetVersionID, req.TargetHistoryID)

	// Retrieve columns generated by database
	var id int64
	var state string
	err = row.Scan(&id, &state)
	if err != nil {
		return entities.Request{}, err
	}

	req.RequestID = id
	req.State = state
	return req, nil
}

func (r PgRequestRepository) SetStatus(request int64, status string) error {

	// store request entity
	stmt, err := r.Db.Prepare("UPDATE request SET state = $2 WHERE id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(request, status)
	if err != nil {
		return err
	}
	return nil
}

func (r PgRequestRepository) GetRequest(request int64) (entities.Request, error) {
	// store request entity
	stmt, err := r.Db.Prepare("SELECT * FROM request WHERE id=$1")
	if err != nil {
		return entities.Request{}, err
	}
	row := stmt.QueryRow(request)

	// Retrieve columns generated by database
	var req entities.Request
	err = row.Scan(&req.RequestID, &req.ArticleID, &req.SourceVersionID, &req.SourceHistoryID, &req.TargetVersionID, &req.TargetHistoryID, &req.State)
	if err != nil {
		return entities.Request{}, err
	}

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
