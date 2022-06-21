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

	// get the latest commit from the git branch after merging
	commit, err := s.Gitrepo.GetLatestCommit(req.ArticleID, req.TargetVersionID)
	if err != nil {
		return err
	}

	// update the commit id of the version in the database
	err = s.Versionrepo.UpdateVersionLatestCommit(req.TargetVersionID, commit)
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

	// get source and target version info
	source, err := s.Versionrepo.GetVersion(req.SourceVersionID)
	if err != nil {
		return models.RequestWithComparison{}, err
	}
	target, err := s.Versionrepo.GetVersion(req.TargetVersionID)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	// ensure that the before-and-after comparison is up to date
	err = s.UpdateRequestComparison(req, source, target)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	// Get the request preview
	before, after, err := s.Gitrepo.GetRequestComparison(req.ArticleID, req.RequestID)
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

// UpdateRequestComparison stores the before and after of the request, if it isn't up-to-date yet
func (s RequestService) UpdateRequestComparison(req entities.Request, source entities.Version, target entities.Version) error {
	if req.Status != "pending" {
		// if not pending anymore, the comparison should not be updated
		return nil
	}

	if req.SourceHistoryID == source.LatestCommitID && req.TargetHistoryID == target.LatestCommitID {
		// if the commits that the request compares are up-to-date, the current comparison can be used
		return nil
	}
	req.SourceHistoryID = source.LatestCommitID
	req.TargetHistoryID = target.LatestCommitID

	// store the preview using git merging and check if there will be conflicts
	conflicted, err := s.Gitrepo.StoreRequestComparison(req)
	if err != nil {
		return err
	}
	req.Conflicted = conflicted

	// now that the comparison has successfully been updated, store the correct commit IDs in the db
	err = s.Repo.UpdateRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (s RequestService) GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
	return s.Repo.GetRequestList(articleId, sourceId, targetId, relatedId)
}
