package interfaces

type ThreadRepository interface {
	CreateThread(aid int64) (int64, error)
}
