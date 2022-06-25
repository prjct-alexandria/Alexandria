package services

import (
	"errors"
	"fmt"
	"mainServer/entities"
	"mainServer/models"
	"mainServer/repositories/interfaces"
	"mainServer/repositories/storer"
	"mainServer/utils"
)

type RequestService struct {
	Repo        interfaces.RequestRepository
	Versionrepo interfaces.VersionRepository
	Storer      storer.Storer
}

func (s RequestService) CreateRequest(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error) {
	// Create request object
	req := entities.Request{
		ArticleID:       article,
		SourceVersionID: sourceVersion,
		TargetVersionID: targetVersion,
	}

	//Check if user is allowed to create request
	isSourceOwner, err := s.Versionrepo.CheckIfOwner(req.SourceVersionID, loggedInAs)
	if err != nil {
		return models.Request{}, errors.New("could not verify version ownership")
	}
	isTargetOwner, err := s.Versionrepo.CheckIfOwner(req.TargetVersionID, loggedInAs)
	if err != nil {
		return models.Request{}, errors.New("could not verify version ownership")
	}
	// TODO make endpoint return 403 Forbidden after this error
	if !(isSourceOwner || isTargetOwner) {
		return models.Request{},
			errors.New(fmt.Sprintf(`request creation forbidden: %v does not own source or target version`, loggedInAs))
	}

	// create the request entity in the db
	req, err = s.Repo.CreateRequest(req)
	if err != nil {
		return models.Request{}, err
	}

	//  directly converts the entity to the model, because they have the exact same fields
	return models.Request(req), nil
}

// RejectRequest rejects the specified request, changing its status
// returns an error if the current user doesn't own the target version
func (s RequestService) RejectRequest(request int64, loggedInAs string) error {

	// get the request information
	req, err := s.Repo.GetRequest(request)
	if err != nil {
		return err
	}

	// get source and target version info for checking owners and latest commits
	source, err := s.Versionrepo.GetVersion(req.SourceVersionID)
	if err != nil {
		return err
	}
	target, err := s.Versionrepo.GetVersion(req.TargetVersionID)
	if err != nil {
		return err
	}

	// check if logged-in user owns the target version
	if !utils.Contains(target.Owners, loggedInAs) {
		return fmt.Errorf("request cannot be rejected, because %v does not own version %v", email, target)
	}

	// update the request comparison one last time before rejecting
	err = s.UpdateRequestComparison(req, source, target)
	if err != nil {
		return err
	}
	// reject the request
	return s.Repo.SetStatus(request, entities.RequestRejected)
}

// AcceptRequest accepts the specified request, changing its status, recording the last commits and committing the merge in git.
// returns an error if the current user doesn't own the target version
func (s RequestService) AcceptRequest(request int64, loggedInAs string) error {

	// get the request information
	req, err := s.Repo.GetRequest(request)
	if err != nil {
		return err
	}
	if req.Conflicted {
		return fmt.Errorf("request %d cannot be accepted, because there would be merge conflicts", request)
	}

	// get source and target version info for checking owners and latest commits
	source, err := s.Versionrepo.GetVersion(req.SourceVersionID)
	if err != nil {
		return err
	}
	target, err := s.Versionrepo.GetVersion(req.TargetVersionID)
	if err != nil {
		return err
	}

	// check if logged-in user owns the target version
	if !utils.Contains(target.Owners, loggedInAs) {
		return fmt.Errorf("request cannot be rejected, because %v does not own version %v", email, target)
	}

	// update the request comparison one last time before accept
	err = s.UpdateRequestComparison(req, source, target)
	if err != nil {
		return err
	}

	// Merge
	commit, err := s.Storer.Merge(req.ArticleID, source.Id, target.Id)
	if err != nil {
		return err
	}

	// Store the commit id in the database
	err = s.Versionrepo.UpdateVersionLatestCommit(target.Id, commit)
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

	// ensure that the before-and-after comparison is up-to-date
	err = s.UpdateRequestComparison(req, source, target)
	if err != nil {
		return models.RequestWithComparison{}, err
	}

	// Get the request preview
	before, after, err := s.Storer.GetRequestComparison(req.ArticleID, req.RequestID)
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
	// if not pending anymore, the comparison should not be updated
	if req.Status != "pending" {
		return nil
	}

	// if the history ID's as in the request entity are the same as in the versions itself, it's up-to-date
	if req.SourceHistoryID == source.LatestCommitID && req.TargetHistoryID == target.LatestCommitID {
		// if the commits that the request compares are up-to-date, the current comparison can be used
		return nil
	}
	req.SourceHistoryID = source.LatestCommitID
	req.TargetHistoryID = target.LatestCommitID

	// store the preview using git merging and check if there will be conflicts
	conflicted, err := s.Storer.StoreRequestComparison(req.ArticleID, req.RequestID, req.SourceVersionID, req.TargetVersionID)
	if err != nil {
		return err
	}

	// now that the comparison has successfully been updated, store the correct commit IDs in the db
	req.Conflicted = conflicted
	err = s.Repo.UpdateRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (s RequestService) GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error) {
	return s.Repo.GetRequestList(articleId, sourceId, targetId, relatedId)
}
