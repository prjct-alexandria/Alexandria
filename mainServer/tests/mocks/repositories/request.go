package repositories

import (
	"mainServer/entities"
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type RequestRepositoryMock struct {
	Mock *mocks.Mock
}

func NewRequestRepositoryMock() RequestRepositoryMock {
	return RequestRepositoryMock{Mock: mocks.NewMock()}
}

var CreateRequestMock func(req entities.Request) (entities.Request, error)

func (m RequestRepositoryMock) CreateRequest(req entities.Request) (entities.Request, error) {
	m.Mock.CallFunc("CreateRequest", &map[string]interface{}{
		"req": req,
	})
	return CreateRequestMock(req)
}

var SetStatusMock func(request int64, status string) error

func (m RequestRepositoryMock) SetStatus(request int64, status string) error {
	m.Mock.CallFunc("SetStatus", &map[string]interface{}{
		"request": request,
		"status":  status,
	})
	return SetStatusMock(request, status)
}

var GetRequestMock func(request int64) (entities.Request, error)

func (m RequestRepositoryMock) GetRequest(request int64) (entities.Request, error) {
	m.Mock.CallFunc("GetRequest", &map[string]interface{}{
		"request": request,
	})
	return GetRequestMock(request)
}

var UpdateRequestMock func(req entities.Request) error

func (m RequestRepositoryMock) UpdateRequest(req entities.Request) error {
	m.Mock.CallFunc("UpdateRequest", &map[string]interface{}{
		"req": req,
	})
	return UpdateRequestMock(req)
}

var GetRequestListMock func(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error)

func (m RequestRepositoryMock) GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
	m.Mock.CallFunc("GetRequestList", &map[string]interface{}{
		"articleId": articleId,
		"sourceId":  sourceId,
		"targetId":  targetId,
		"relatedId": relatedId,
	})
	return GetRequestListMock(articleId, sourceId, targetId, relatedId)
}
