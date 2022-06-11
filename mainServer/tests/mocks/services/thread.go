package services

import "mainServer/models"

// ThreadServiceMock mocks class using publicly modifiable mock functions
type ThreadServiceMock struct {
	// mock tracks what functions were called and with what parameters
	Called *map[string]bool
	Params *map[string]map[string]interface{}
}

// NewThreadServiceMock initializes a mock with variables that are passed by reference,
// so the values can be retrieved from anywhere in the program
func NewThreadServiceMock() ThreadServiceMock {
	return ThreadServiceMock{
		Called: &map[string]bool{},
		Params: &map[string]map[string]interface{}{},
	}
}

func (m ThreadServiceMock) StartThread(thread models.Thread, aid int64, sid int64) (int64, error) {
	panic("implement me")
}
