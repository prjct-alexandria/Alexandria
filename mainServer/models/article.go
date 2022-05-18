package models

type Article struct {
	Id     string `json:"id"`
	Title  string `json:"title"`
	Owners []User `json:"owners"`
}

type ArticleCreationForm struct {
	Title  string   `json:"title" binding:"required"`
	Owners []string `json:"owners" binding:"required"`
}
