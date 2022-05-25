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
	SpecificId int64
}
