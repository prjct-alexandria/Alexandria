package entities

type Request struct {
	RequestID       int64
	ArticleID       int64
	SourceVersionID int64
	SourceHistoryID int64
	TargetVersionID int64
	TargetHistoryID int64
	State           string
}

const (
	RequestAccepted string = "accepted"
	RequestRejected        = "rejected"
	RequestPending         = "pending"
)
