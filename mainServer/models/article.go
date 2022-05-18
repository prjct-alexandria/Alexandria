package models

type ArticleCreationForm struct {
	Title  string   `json:"title" binding:"required"`
	Owners []string `json:"owners" binding:"required"`
}
