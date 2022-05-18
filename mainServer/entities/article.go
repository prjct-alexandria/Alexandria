package entities

type Article struct {
	Id     string   `json:"id" binding:"required"`
	Title  string   `json:"title" binding:"required"`
	Owners []string `json:"owners" binding:"required"`
}
