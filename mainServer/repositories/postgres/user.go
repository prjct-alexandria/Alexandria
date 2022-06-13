package postgres

import (
	"database/sql"
	"fmt"
	"mainServer/entities"
)

type PgUserRepository struct {
	Db *sql.DB
}

func NewPgUserRepository(db *sql.DB) PgUserRepository {
	repo := PgUserRepository{db}
	err := repo.createUserTable()
	if err != nil {
		panic(err)
	}
	return repo
}

func (r PgUserRepository) CreateUser(user entities.User) error {
	stmt, err := r.Db.Prepare("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}

	_, err = stmt.Exec(user.Name, user.Email, user.Pwd)
	return err
}

func (r PgUserRepository) GetFullUserByEmail(email string) (entities.User, error) {
	var user entities.User
	stmt, err := r.Db.Prepare("SELECT * FROM users WHERE Email = $1")

	if err != nil {
		return user, err
	}

	row := stmt.QueryRow(email)
	if err := row.Scan(&user.Name, &user.Email, &user.Pwd); err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("GetFullUserByEmail %s: no such user", email)
		}
		return user, fmt.Errorf("GetFullUserByEmail %s: %v", email, err)
	}
	return user, nil
}

func (r PgUserRepository) UpdateUser(email string, user entities.User) error {
	//To be implemented
	return nil
}

func (r PgUserRepository) DeleteUser(email string) error {
	//To be implemented
	return nil
}

func (r PgUserRepository) createUserTable() error {
	//_, err := r.Db.Exec("DROP TABLE IF EXISTS users")

	_, err := r.Db.Exec(`CREATE TABLE IF NOT EXISTS users (
    	username varchar(64) NOT NULL,
    	email varchar(64) NOT NULL,
    	password varchar(64) NOT NULL,
    	PRIMARY KEY (email)
    )`)
	return err
}
