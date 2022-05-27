package models

type Version struct {
	ArticleID int64    `json:"articleID"`
	Id        int64    `json:"versionID"`
	Title     string   `json:"title"`
	Owners    []string `json:"owners"`
	Content   string   `json:"content"`
}

type VersionCreationForm struct {
	SourceVersionID int64    `json:"sourceVersionID" binding:"required"`
	Title           string   `json:"title" binding:"required"`
	Owners          []string `json:"owners" binding:"required"`
}
