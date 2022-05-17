package repositories

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
		return PgUserRepository{}
	}
	return repo
}

func (r PgUserRepository) CreateUser(user entities.User) error {
	_, err := r.Db.Exec("INSERT INTO users (username, email, password) VALUES ('" +
		user.Name + "', '" +
		user.Email + "', '" +
		user.Pwd + "')")
	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}
	return err
}

func (r PgUserRepository) GetFullUserByEmail(email string) (entities.User, error) {
	var user entities.User

	row := r.Db.QueryRow("SELECT * FROM users WHERE Email = '" + email + "'")
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
