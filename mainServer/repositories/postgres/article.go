package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
)

type PgArticleRepository struct {
	Db *sql.DB
}

func NewPgArticleRepository(db *sql.DB) PgArticleRepository {
	repo := PgArticleRepository{db}

	err := repo.createArticleTable()
	if err != nil {
		panic(err)
	}

	err = repo.createArticleVersionsTable()
	if err != nil {
		panic(err)
	}
	return repo
}

func (r PgArticleRepository) CreateArticle() (entities.Article, error) {
	result, err := r.Db.Exec("INSERT INTO article (id) VALUES (DEFAULT)")
	if err != nil {
		return entities.Article{}, err
	}

	// read the generated id to return, should always be here if the query itself returned no error
	id, _ := result.LastInsertId()
	article := entities.Article{Id: id}

	return article, nil
}

func (r PgArticleRepository) LinkVersion(articleID int64, versionID int64) error {
	_, err := r.Db.Exec(fmt.Sprintf("INSERT INTO articleVersions (versionID, articleID) VALUES (%d, %d)", versionID, articleID))
	return err
}

func (r PgArticleRepository) createArticleTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS article (
    			id SERIAL PRIMARY KEY
    )`)
	return err
}

func (r PgArticleRepository) createArticleVersionsTable() error {
	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS articleVersions (
    			versionID int PRIMARY KEY,
    			articleID int NOT NULL,
    			FOREIGN KEY (articleID) REFERENCES article(id),
    			FOREIGN KEY (versionID) REFERENCES version(id)
    )`)
	return err
}
