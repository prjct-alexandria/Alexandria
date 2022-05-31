package entities

type Request struct {
	RequestID       int64
	ArticleID       int64
	SourceVersionID int64
	SourceHistoryID string
	TargetVersionID int64
	TargetHistoryID string
	State           string
}

const (
	RequestAccepted string = "accepted"
	RequestRejected        = "rejected"
	RequestPending         = "pending"
)
