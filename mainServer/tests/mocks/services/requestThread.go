package services

import "mainServer/models"

var GetRequestThreadsMock func(aid int64, rid int64) ([]models.Thread, error)

// RequestThreadServiceMock mocks class using publicly modifiable mock functions
type RequestThreadServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

// NewRequestThreadServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewRequestThreadServiceMock() RequestThreadServiceMock {
	return RequestThreadServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

func (m RequestThreadServiceMock) StartRequestThread(rid int64, tid int64, loggedInAs string) (int64, error) {
	panic("implement me")
}

func (m RequestThreadServiceMock) GetRequestThreads(aid int64, rid int64) ([]models.Thread, error) {
	(*m.Called)["GetRequestThreads"] = true
	(*m.Params)["GetRequestThreads"] = map[string]interface{}{
		"aid": aid,
		"rid": rid,
	}
	return GetRequestThreadsMock(aid, rid)
}
