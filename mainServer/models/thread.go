package models

import (
	"mainServer/entities"
)

type Thread struct {
	Id         int64
	ArticleId  int64
	SpecificId int64
	Comment    []entities.Comment
}

type ReturnThreadIds struct {
	Id        int64
	ThreadId  int64
	CommentId int64
}
