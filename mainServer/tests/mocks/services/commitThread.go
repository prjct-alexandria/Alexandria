package services

import "mainServer/models"

var GetCommitThreadsMock func(aid int64, cid int64) ([]models.Thread, error)

// CommitThreadServiceMock mocks class using publicly modifiable mock functions
type CommitThreadServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

// NewCommitThreadServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewCommitThreadServiceMock() CommitThreadServiceMock {
	return CommitThreadServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

func (m CommitThreadServiceMock) StartCommitThread(thread models.Thread, tid int64) (int64, error) {
	panic("implement me")
}

func (m CommitThreadServiceMock) GetCommitThreads(aid int64, cid int64) ([]models.Thread, error) {
	(*m.Called)["GetCommitThreads"] = true
	(*m.Params)["GetCommitThreads"] = map[string]interface{}{
		"aid": aid,
		"cid": cid,
	}
	return GetCommitThreadsMock(aid, cid)
}
