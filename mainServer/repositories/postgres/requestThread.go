package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
	"strconv"
)

type PgRequestThreadRepository struct {
	Db *sql.DB
}

func NewPgRequestThreadRepository(db *sql.DB) PgRequestThreadRepository {
	repo := PgRequestThreadRepository{db}
	err := repo.createRequestThreadTable()
	if err != nil {
		return PgRequestThreadRepository{}
	}
	return repo
}

func (r PgRequestThreadRepository) createRequestThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS RequestThread (
    	RequestThreadId SERIAL,
    	RequestId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (RequestThreadId)
    )`)
	return err
}

func (r PgRequestThreadRepository) CreateRequestThread(thread models.ThreadNoId, tid string) (int64, error) {
	row := r.Db.QueryRow("INSERT INTO RequestThread (requestId, threadId) " +
		"VALUES ('" +
		strconv.FormatInt(thread.SpecificId, 10) + "', '" +
		tid + "')" +
		"RETURNING RequestThreadId")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateThread: %v", err)
	}

	return id, err
}
