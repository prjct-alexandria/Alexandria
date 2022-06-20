package storer

import "mainServer/repositories/storer/git"

// GetVersionFiles creates a .zip with all files of an article version at the returned path
func (s *Storer) GetVersionFiles(article int64, version int64, filename string) (string, error) {
	s.pool.Lock(article)
	defer s.pool.Unlock(article)

	// Create GitRepo that refers to article
	path, err := s.fs.GetArticlePath(article)
	if err != nil {
		return "", err
	}
	repo := git.NewRepo(path)

	// Checkout version
	err = repo.CheckoutBranch(version)
	if err != nil {
		return "", err
	}

	// Add a .zip with the files to a cache folder and return the path
	path, err = s.fs.MakeDownloadZip(filename)
	if err != nil {
		return "", nil
	}
	return path, nil
}
