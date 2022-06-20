package storer

import (
	"mainServer/repositories/storer/filesystem"
	"mainServer/repositories/storer/locking"
	"mainServer/server/config"
	"mainServer/utils/clock"
)

// Storer handles every functionality that accesses (git) files in the OS filesystem
// Manages thread safety and should not be bypassed.
type Storer struct {
	pool  *locking.MutexPool
	fs    *filesystem.FileSystem
	clock clock.Clock
}

// NewStorer returns Storer, initialized with the config. Creates necessary paths
func NewStorer(cfg *config.StorerConfig) Storer {
	return Storer{
		pool:  locking.NewMutexPool(cfg.MutexCount),
		fs:    filesystem.NewFileSystem(cfg.Path),
		clock: clock.RealClock{},
	}
}

// SetClock changes the clock that is used to set creation dates of commits
// Can be used to supply a mock for testing purposes
func (s Storer) SetClock(clock clock.Clock) {
	s.clock = clock
}
