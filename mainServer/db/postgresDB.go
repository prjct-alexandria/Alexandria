package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"mainServer/server/config"
)

func Connect(cfg *config.DatabaseConfig) *sql.DB {
	// fill in connection string
	var sslMode string
	if cfg.Url.UseSSL {
		sslMode = "enable"
	} else {
		sslMode = "disable"
	}
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Url.Host, cfg.Url.Port, cfg.User, cfg.Pwd, cfg.DBName, sslMode)

	// open database
	db, err := sql.Open("postgres", conn)
	CheckError(err)

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
