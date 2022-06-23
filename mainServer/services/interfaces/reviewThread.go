package interfaces

type ReviewThreadService interface {
	StartReviewThread(rid int64, tid int64) (int64, error)
}
