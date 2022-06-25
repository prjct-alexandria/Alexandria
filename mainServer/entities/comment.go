package entities

type Comment struct {
	Id           int64  `json:"id"`
	AuthorId     string `json:"authorId"        binding:"required"`
	ThreadId     int64  `json:"threadId"`
	Content      string `json:"content"         binding:"required"`
	CreationDate string `json:"creationDate"    binding:"required"`
}
