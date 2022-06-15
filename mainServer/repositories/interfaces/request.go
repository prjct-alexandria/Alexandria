package interfaces

import (
	"mainServer/entities"
	"mainServer/models"
)

type RequestRepository interface {

	// CreateRequest creates the specified entity in the database
	// request ID is generated and attached in the returned entity
	CreateRequest(req entities.Request) (entities.Request, error)

	// SetStatus sets the status of the specified request
	SetStatus(request int64, status string) error

	// GetRequest returns the request entity with the specified ID
	GetRequest(request int64) (entities.Request, error)

	// UpdateRequest replaces all the fields with the entity specified, in the row with matching id
	UpdateRequest(req entities.Request) error

	// GetRequestList returns a list of request models related to specified article, with possible filters
	GetRequestList(articleId int64, sourceId int64, targetId int64, relatedId int64) ([]models.RequestListElement, error)
}
