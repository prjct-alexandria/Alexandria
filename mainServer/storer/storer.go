package storer

import (
	"mainServer/repositories/interfaces"
	"mainServer/server/config"
	"mainServer/utils/clock"
	"os"
	"path/filepath"
)

// Storer handles every functionality that accesses (git) files in the OS filesystem
// Manages thread safety and should not be bypassed.
type Storer struct {
	path        string
	pool        *MutexPool
	clock       clock.Clock
	versionRepo interfaces.VersionRepository
	download    DownloadStorer
	git         GitStorer
}

// NewStorer returns Storer, initialized with the config. Creates necessary paths
func NewStorer(version interfaces.VersionRepository, cfg *config.FileSystemConfig) Storer {

	// make folder for git repos
	err := os.MkdirAll(filepath.Join(cfg.Path, "persistent", "articles"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// make folder for request before-and-afters
	err = os.MkdirAll(filepath.Join(cfg.Path, "persistent", "requests"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// make (cache) folder for download zip's
	err = os.MkdirAll(filepath.Join(cfg.Path, "cache", "download"), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return Storer{
		path:        cfg.Path,
		pool:        NewMutexPool(cfg.MutexCount),
		clock:       clock.RealClock{},
		versionRepo: version,
	}
}

// GetVersionFiles creates a .zip with all files of an article version at the returned path
func (s *Storer) GetVersionFiles(article int64, version int64, filename string) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Checkout version
	err := s.git.CheckoutBranch(article, version)
	if err != nil {
		return "", err
	}

	path, err := s.download.MakeDownloadZip(filename)
	if err != nil {
		return "", nil
	}
	return path, nil
}
