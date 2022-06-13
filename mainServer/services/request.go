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
	req, err = s.UpdateRequestPreview(req)
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

	// check if logged-in user owns this version
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
	if req.Conflicted {
		return fmt.Errorf("request %d cannot be accepted, because there would be merge conflicts", request)
	}

	// check if logged-in user owns this version
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

// GetRequest returns a request, including the before and after versions
func (s RequestService) GetRequest(request int64) (models.RequestWithComparison, error) {
	// get request info
	req, err := s.Repo.GetRequest(request)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	if req.Status == "pending" {
		// If the request is still active/pending,
		// make sure that it's up-to-date with the latest commit on both versions
		req, err = s.UpdateRequestPreview(req)
		if err != nil {
			return models.RequestWithComparison{}, err
		}
	}
	before, after, err := s.Gitrepo.GetRequestPreview(req.ArticleID, req.SourceHistoryID, req.TargetHistoryID)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	// get both request versions
	source, err := s.Versionrepo.GetVersion(req.SourceVersionID)
	if err != nil {
		return models.RequestWithComparison{}, err
	}
	target, err := s.Versionrepo.GetVersion(req.TargetVersionID)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	// insert before and after contents in the version models
	return models.RequestWithComparison{
		Request: models.Request(req),
		Source: models.Version{
			ArticleID: source.ArticleID,
			Id:        source.Id,
			Title:     source.Title,
			Owners:    source.Owners,
			Status:    source.Status,
		},
		Target: models.Version{
			ArticleID: target.ArticleID,
			Id:        target.Id,
			Title:     target.Title,
			Owners:    target.Owners,
			Status:    target.Status,
		},
		Before: before,
		After:  after,
	}, nil
}

// UpdateRequestPreview stores the before and after of the article in a cache, if it does not exist yet
// updates the request in the db with the latest commits and the conflict status, also returns entity
func (s RequestService) UpdateRequestPreview(req entities.Request) (entities.Request, error) {
	var err error // declare error in advance, so multiple assignment of req fields and err works

	// record the current most recent history/commit IDs of both versions (branches)
	req.SourceHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.SourceVersionID)
	if err != nil {
		return entities.Request{}, err
	}
	req.TargetHistoryID, err = s.Gitrepo.GetLatestCommit(req.ArticleID, req.TargetVersionID)
	if err != nil {
		return entities.Request{}, err
	}

	// don't do anything if it already exists
	exists, err := s.Gitrepo.RequestPreviewExists(req)
	if err != nil {
		return entities.Request{}, err
	}
	if exists {
		return req, nil
	}

	// store the preview using git merging and check if there will be conflicts
	success, err := s.Gitrepo.StoreRequestPreview(req)
	if err != nil {
		return entities.Request{}, err
	}
	req.Conflicted = !success

	// store the updated request in the db
	err = s.Repo.UpdateRequest(req)
	if err != nil {
		return entities.Request{}, err
	}

	return req, nil
}
