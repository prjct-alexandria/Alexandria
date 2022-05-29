package models

type ArticleCreationForm struct {
	Title  string   `json:"title" binding:"required"`
	Owners []string `json:"owners" binding:"required"`
}

type ArticleListElement struct {
	Id     int64    `json:"articleId" binding:"required"`
	Title  string   `json:"title" binding:"required"`
	Owners []string `json:"owners" binding:"required"`
	//Creation date?
}
