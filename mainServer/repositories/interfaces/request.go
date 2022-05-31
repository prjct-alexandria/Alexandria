package interfaces

import "mainServer/entities"

type RequestRepository interface {

	// CreateRequest creates the specified entity in the database
	// request ID is generated and attached in the returned entity
	CreateRequest(req entities.Request) (entities.Request, error)
}
