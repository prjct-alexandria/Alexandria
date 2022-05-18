package interfaces

import "mainServer/entities"

type VersionRepository interface {
	CreateVersion(version entities.Version) error
}
