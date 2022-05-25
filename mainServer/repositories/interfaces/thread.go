package interfaces

type ThreadRepository interface {
	CreateThread(aid string) (int64, error)
}
