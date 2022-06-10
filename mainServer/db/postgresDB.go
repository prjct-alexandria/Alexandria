package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "AlexandriaPG"
)

func Connect() *sql.DB {
	// connection string
	// TODO: use config
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	//defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected to database!")
	return db
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
