package models

type Request struct {
	RequestID       int64  `json:"requestID"`
	ArticleID       int64  `json:"articleID"`
	SourceVersionID int64  `json:"sourceVersionID"`
	SourceHistoryID string `json:"sourceHistoryID"`
	TargetVersionID int64  `json:"targetVersionID"`
	TargetHistoryID string `json:"targetHistoryID"`
	State           string `json:"state"`
}

type RequestCreationForm struct {
	SourceVersionID int64 `json:"sourceVersionID" binding:"required"`
	TargetVersionID int64 `json:"targetVersionID" binding:"required"`
}
