package postgres

import (
	"database/sql"
	"mainServer/entities"
)

type PgArticleRepository struct {
	Db *sql.DB
}

func (r PgArticleRepository) UpdateMainVersion(article int64, version int64) error {
	stmt, err := r.Db.Prepare(`UPDATE article SET mainVersionID=$2 WHERE id=$1`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(article, version)
	return err
}

func NewPgArticleRepository(db *sql.DB) PgArticleRepository {
	repo := PgArticleRepository{db}

	err := repo.createArticleTable()
	if err != nil {
		panic(err)
	}
	return repo
}

func (r PgArticleRepository) CreateArticle() (entities.Article, error) {
	row := r.Db.QueryRow("INSERT INTO article (id) VALUES (DEFAULT) RETURNING id")

	// Because QueryRow instead of Exec was used, with the RETURNING keyword,
	// The generated id can be retrieved
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return entities.Article{}, err
	}

	article := entities.Article{Id: id}
	return article, nil
}

func (r PgArticleRepository) createArticleTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS article (
    			id SERIAL PRIMARY KEY,
            	mainVersionID int
    )`)
	return err
}
