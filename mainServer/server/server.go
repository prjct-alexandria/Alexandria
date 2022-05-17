package server

import "mainServer/db"

func Init() {
	Database := db.Connect()
	router := SetUpRouter(Database)

	router.Run("localhost:8080")
}
