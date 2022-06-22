package models

import (
	"mainServer/entities"
)

type Thread struct {
	Id         int64              `json:"id"`
	ArticleId  int64              `json:"articleId"       binding:"required"`
	SpecificId string             `json:"specificId"`
	Comments   []entities.Comment `json:"comments"         binding:"required"`
}

type ReturnThreadIds struct {
	Id        int64 `json:"id"              binding:"required"`
	ThreadId  int64 `json:"threadId"        binding:"required"`
	CommentId int64 `json:"CommentId"       binding:"required"`
}
