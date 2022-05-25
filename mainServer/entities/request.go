package entities

type Request struct {
	articleID       int64
	requestID       int64
	sourceVersionID int64
	sourceHistoryID int64
	targetVersionID int64
	targetHistoryID int64
	state           string
}

const (
	RequestAccepted string = "accepted"
	RequestRejected        = "rejected"
	RequestPending         = "pending"
)
