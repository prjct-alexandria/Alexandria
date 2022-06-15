package interfaces

type ReviewThreadRepository interface {
	CreateReviewThread(rid int64, tid int64) (int64, error)
}
