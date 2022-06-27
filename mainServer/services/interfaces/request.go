package interfaces

import (
	"mainServer/entities"
	"mainServer/models"
)

type RequestService interface {
	CreateRequest(article int64, sourceVersion int64, targetVersion int64, loggedInAs string) (models.Request, error)

	// RejectRequest rejects the specified request, changing its status
	// returns an error if the current user doesn't own the target version
	RejectRequest(request int64, loggedInAs string) error

	// AcceptRequest accepts the specified request, changing its status, recording the last commits and committing the merge in git.
	// returns an error if the current user doesn't own the target version
	AcceptRequest(request int64, loggedInAs string) error

	// GetRequest returns a request, including the before and after versions
	GetRequest(request int64) (models.RequestWithComparison, error)

	// UpdateRequestComparison stores the before and after of the request, if it isn't up-to-date yet
	UpdateRequestComparison(req entities.Request, source entities.Version, target entities.Version) error

	GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error)
}
