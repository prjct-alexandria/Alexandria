package entities

type Comment struct {
	Id           int64
	AuthorId     string
	ThreadId     int64
	Content      string
	CreationDate string
}
