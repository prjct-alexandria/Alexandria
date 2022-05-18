package models

type Version struct {
	ArticleID string   `json:"articleID"`
	Id        string   `json:"versionID"`
	Title     string   `json:"title"`
	Owners    []string `json:"owners"`
}
