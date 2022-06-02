package interfaces

import "mainServer/entities"

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
}
