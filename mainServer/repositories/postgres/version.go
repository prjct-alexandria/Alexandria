package postgres

import (
	"database/sql"
	"mainServer/entities"
)

type PgVersionRepository struct {
	Db *sql.DB
}

func NewPgVersionRepository(db *sql.DB) PgVersionRepository {
	repo := PgVersionRepository{db}
	//err := repo.createUserTable()
	//if err != nil {
	//	return PgUserRepository{}
	//}
	return repo
}

func (repo PgVersionRepository) CreateVersion(version entities.Version) error {
	
}
