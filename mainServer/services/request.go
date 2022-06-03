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

	// create the request entity in the db
	req, err := s.Repo.CreateRequest(req)
	if err != nil {
		return models.Request{}, err
	}

	// create the request preview cache, also updates the entity in the db for whether it has conflicts
	err = s.UpdateRequestPreview(req)
	if err != nil {
		return models.Request{}, err
	}

	//  directly converts the entity to the model, because they have the exact same fields
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

// AcceptRequest accepts the specified request, changing its status, recording the last commits and committing the merge in git.
// returns an error if the current user doesn't own the target version
func (s RequestService) AcceptRequest(request int64, email string) error {

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
		return fmt.Errorf("request cannot be accepted, because %v does not own version %v", email, target)
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

	// commit the merge in git
	err = s.Gitrepo.Merge(req.ArticleID, req.SourceVersionID, req.TargetVersionID)
	if err != nil {
		return err
	}

	// accept the request
	return s.Repo.SetStatus(request, entities.RequestAccepted)
}

// UpdateRequestPreview stores the before and after of the article in a cache,
// updates the request in the db with the latest commits and the conflict status
func (s RequestService) UpdateRequestPreview(req entities.Request) error {
	var err error // declare error in advance, so multiple assignment of req fields and err works

	// record the current most recent history/commit IDs of both versions (branches)
	req.SourceHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.SourceVersionID)
	if err != nil {
		return err
	}
	req.TargetHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.TargetVersionID)
	if err != nil {
		return err
	}

	// store the preview using git merging and check if there will be conflicts
	success, err := s.Gitrepo.StoreRequestPreview(req)
	if err != nil {
		return err
	}
	req.Conflicted = !success

	// store the updated request in the db
	err = s.Repo.UpdateRequest(req)
	if err != nil {
		return err
	}

	return nil
}
