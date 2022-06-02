package services

import (
	"fmt"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories"
	"mainServer/repositories/interfaces"
)

type RequestService struct {
	Repo        interfaces.RequestRepository
	Versionrepo interfaces.VersionRepository
	Gitrepo     repositories.GitRepository
}

func (s RequestService) CreateRequest(article int64, sourceVersion int64, targetVersion int64) (models.Request, error) {
	req := entities.Request{
		ArticleID:       article,
		SourceVersionID: sourceVersion,
		TargetVersionID: targetVersion,
	}

	req, err := s.Repo.CreateRequest(req)
	if err != nil {
		return models.Request{}, err
	}

	//  directly converts the entity to the model, because they have the exact same fiels
	return models.Request(req), nil
}

// RejectRequest rejects the specified request, changing its status
// returns an error if the current user doesn't own the target version
func (s RequestService) RejectRequest(request int64, email string) error {

	// get the request information
	req, err := s.Repo.GetRequest(request)
	if err != nil {
		return err
	}

	// check who the owner is
	target := req.TargetVersionID
	ok, err := s.Versionrepo.CheckIfOwner(target, email)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("request cannot be rejected, because %v does not own version %v", email, target)
	}

	// record the current most recent history/commit IDs of both versions (branches)
	req.SourceHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.SourceVersionID)
	if err != nil {
		return err
	}
	req.TargetHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.TargetVersionID)
	if err != nil {
		return err
	}
	err = s.Repo.UpdateRequest(req)
	if err != nil {
		return err
	}

	// reject the request
	return s.Repo.SetStatus(request, entities.RequestRejected)
}
