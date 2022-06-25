package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
	"strconv"
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

func (r PgArticleRepository) GetAllArticles() ([]entities.Article, error) {
	stmt, err := r.Db.Prepare(`SELECT id, mainversionid FROM article`)

	var list []entities.Article
	rows, err := stmt.Query()
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var entity entities.Article
		err := rows.Scan(&entity.Id, &entity.MainVersionID)
		if err != nil {
			fmt.Printf("GetAllArticles: %v\n", err.Error())
			continue
		}
		list = append(list, entity)
	}

	return list, err
}

func (r PgArticleRepository) GetMainVersion(aid int64) (int64, error) {
	var mvid int64

	row := r.Db.QueryRow("SELECT mainversionid FROM article WHERE id=" + strconv.FormatInt(aid, 10))
	if err := row.Scan(&mvid); err != nil {
		if err == sql.ErrNoRows {
			return mvid, fmt.Errorf("GetMainVersion no such article")
		}
		return mvid, fmt.Errorf("GetMainVersion")
	}
	return mvid, nil
}
