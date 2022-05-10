package server

func Init() {
	router := SetUpRouter()
	router.Run("localhost:8080")
}
