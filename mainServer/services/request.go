package services

import (
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
)

type RequestService struct {
	Repo interfaces.RequestRepository
}

func (s RequestService) CreateRequest(article int64, sourceVersion int64, targetVersion int64, sourceHistory string, targetHistory string) (models.Request, error) {
	req := entities.Request{
		ArticleID:       article,
		SourceVersionID: sourceVersion,
		SourceHistoryID: sourceHistory,
		TargetVersionID: targetVersion,
		TargetHistoryID: targetHistory,
	}

	req, err := s.Repo.CreateRequest(req)
	if err != nil {
		return models.Request{}, err
	}

	//  directly converts the entity to the model, because they have the exact same fiels
	return models.Request(req), nil
}
