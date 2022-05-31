package models

type Version struct {
	ArticleID int64    `json:"articleID"`
	Id        int64    `json:"versionID"`
	Title     string   `json:"title"`
	Owners    []string `json:"owners"`
	Content   string   `json:"content"`
	Status    string   `json:"status"`
}
