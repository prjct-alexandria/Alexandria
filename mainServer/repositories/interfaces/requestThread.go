package interfaces

type RequestThreadRepository interface {
	CreateRequestThread(rid int64, tid int64) (int64, error)
}
