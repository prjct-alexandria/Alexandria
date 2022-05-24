package models

import "mainServer/entities"

type CommitThread struct {
	Id        int64
	ArticleId int64
	CommitId  int64
	ThreadId  int64
	Comment   []entities.Comment
}

type CommitThreadNoId struct {
	ArticleId int64
	CommitId  int64
	ThreadId  int64
	Comment   []entities.Comment
}
