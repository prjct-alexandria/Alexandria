package repositories

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FilesystemRepository struct {
	Path string
}

func NewFilesystemRepository(path string) (FilesystemRepository, error) {
	repo := FilesystemRepository{path}

	err := os.MkdirAll(filepath.Join(repo.Path, "cache", "downloads"), os.ModePerm)
	return repo, err
}

func (repo FilesystemRepository) AddFilesInDirToZip(zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Wrapped in function to allow for "defer file.close()"
		err := func() error {
			if file.IsDir() {
				//Check if it is not the git folder
				if file.Name() != ".git" {
					err := repo.AddFilesInDirToZip(zipWriter, filepath.Join(dirPath, file.Name()), file.Name())
					if err != nil {
						return err
					}
				}
			} else {
				//Create file in zip
				zipFile, err := zipWriter.Create(filepath.Join(dirInZip, file.Name()))
				if err != nil {
					return err
				}

				//Open file on branch
				fileReader, err := os.Open(filepath.Join(dirPath, file.Name()))
				if err != nil {
					return err
				}
				defer fileReader.Close()

				//Copy contents over
				_, err = io.Copy(zipFile, fileReader)
				if err != nil {
					return err
				}
			}
			return nil
		}()

		if err != nil {
			return err
		}
	}
	return nil
}
