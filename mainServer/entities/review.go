package entities

type Review struct {
	Id        int64
	ArticleId int64
	CommitId  int64
	AuthorId  string
}
