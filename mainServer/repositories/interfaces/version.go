package interfaces

import "mainServer/entities"

type VersionRepository interface {

	// CreateVersion makes a new version from the specified entity.
	// Will ignore the specified ID and generate a new one.
	// Created entity is returned.
	CreateVersion(version entities.Version) (entities.Version, error)

	// GetVersion returns version entity from database
	GetVersion(version int64) (entities.Version, error)

	// GetVersionsByArticle gets the version entities related to a specific article, links the owners
	GetVersionsByArticle(article int64) ([]entities.Version, error)
}
