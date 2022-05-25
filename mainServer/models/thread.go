package models

import "mainServer/entities"

type Thread struct {
	Id         int64
	ArticleId  int64
	CommitId   int64
	SpecificId int64
	Comment    []entities.Comment
}

type ThreadNoId struct {
	ArticleId  int64
	CommitId   int64
	SpecificId int64
	Comment    []entities.Comment
}

type ReturnIds struct {
	ThreadId   int64
	CommentId  []int64
	SpecificId int64
}

type CommentNoId struct {
	AuthorId     string
	ThreadId     int64
	Content      string
	CreationDate string
}
