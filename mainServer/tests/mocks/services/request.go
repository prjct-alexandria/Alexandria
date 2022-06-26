package services

import (
	"mainServer/entities"
	"mainServer/models"
	mocks "mainServer/tests/util"
)

type RequestServiceMock struct {
	Mock *mocks.Mock
}

func NewRequestServiceMock() RequestServiceMock {
	return RequestServiceMock{Mock: mocks.NewMock()}
}

var CreateRequestMock func(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error)

func (m RequestServiceMock) CreateRequest(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
	m.Mock.CallFunc("CreateRequest", &map[string]interface{}{
		"article":       article,
		"sourceVersion": sourceVersion,
		"targetVersion": targetVersion,
		"loggedInAs":    loggedInAs,
	})
	return CreateRequestMock(article, sourceVersion, targetVersion, loggedInAs)
}

var RejectRequestMock func(request int64, loggedInAs string) error

func (m RequestServiceMock) RejectRequest(request int64, loggedInAs string) error {
	m.Mock.CallFunc("RejectRequest", &map[string]interface{}{
		"request":    request,
		"loggedInAs": loggedInAs,
	})
	return RejectRequestMock(request, loggedInAs)
}

var AcceptRequestMock func(request int64, loggedInAs string) error

func (m RequestServiceMock) AcceptRequest(request int64, loggedInAs string) error {
	m.Mock.CallFunc("AcceptRequest", &map[string]interface{}{
		"request":    request,
		"loggedInAs": loggedInAs,
	})
	return AcceptRequestMock(request, loggedInAs)
}

var GetRequestMock func(request int64) (models.RequestWithComparison, error)

func (m RequestServiceMock) GetRequest(request int64) (models.RequestWithComparison, error) {
	m.Mock.CallFunc("GetRequest", &map[string]interface{}{
		"request": request,
	})
	return GetRequestMock(request)
}

var UpdateRequestComparisonMock func(req entities.Request, source entities.Version, target entities.Version) error

func (m RequestServiceMock) UpdateRequestComparison(req entities.Request, source entities.Version, target entities.Version) error {
	m.Mock.CallFunc("UpdateRequestComparison", &map[string]interface{}{
		"req":    req,
		"source": source,
		"target": target,
	})
	return UpdateRequestComparisonMock(req, source, target)
}

var GetRequestListMock func(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error)

func (m RequestServiceMock) GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
	m.Mock.CallFunc("GetRequestList", &map[string]interface{}{
		"articleId": articleId,
		"sourceId":  sourceId,
		"targetId":  targetId,
		"relatedId": relatedId,
	})
	return GetRequestListMock(articleId, sourceId, targetId, relatedId)
}
