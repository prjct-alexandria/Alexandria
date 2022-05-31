package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/models"
	"strconv"
)

type PgReviewThreadRepository struct {
	Db *sql.DB
}

func NewPgReviewThreadRepository(db *sql.DB) PgReviewThreadRepository {
	repo := PgReviewThreadRepository{db}
	err := repo.createReviewThreadTable()
	if err != nil {
		return PgReviewThreadRepository{}
	}
	return repo
}

func (r PgReviewThreadRepository) createReviewThreadTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS ReviewThread (
    	ReviewThreadId SERIAL,
    	ReviewId BIGINT NOT NULL,
        threadId BIGINT NOT NULL,
    	PRIMARY KEY (ReviewThreadId)
    )`)
	return err
}

func (r PgReviewThreadRepository) CreateReviewThread(thread models.Thread, tid string) (int64, error) {
	row := r.Db.QueryRow("INSERT INTO ReviewThread (reviewId, threadId) " +
		"VALUES ('" +
		strconv.FormatInt(thread.SpecificId, 10) + "', '" +
		tid + "')" +
		"RETURNING ReviewThreadId")
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreateThread: %v", err)
	}

	return id, err
}
