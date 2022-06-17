package storer

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"mainServer/server/config"
	"os"
	"path/filepath"
)

type DownloadStorer struct {
	Path string
}

func NewFilesystemRepository(cfg *config.FileSystemConfig) DownloadStorer {
	path := filepath.Join(cfg.Path, "cache", "downloads")
	repo := DownloadStorer{Path: path}

	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return repo
}

// MakeDownloadZip
func (s DownloadStorer) MakeDownloadZip(filename string) (string, error) {
	// TODO: check if this can lead to some sort of injection
	path := filepath.Join(s.Path, "cache", "downloads", filename+".zip")

	// Create the (empty) zip file
	versionZip, err := os.Create(path)
	defer versionZip.Close()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Create a writer for the file, with the zip library
	zipWriter := zip.NewWriter(versionZip)
	defer zipWriter.Close()

	err = s.addFilesInDirToZip(zipWriter, path, "")
	if err != nil {
		return "", err
	}

}

func (s DownloadStorer) addFilesInDirToZip(zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Wrapped in function to allow for "defer file.close()"
		err := s.addFileInDirToZip(file, zipWriter, dirPath, dirInZip)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s DownloadStorer) addFileInDirToZip(file fs.FileInfo, zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	if file.IsDir() {
		//Check if it is not the git folder
		if file.Name() != ".git" {
			err := s.addFilesInDirToZip(zipWriter, filepath.Join(dirPath, file.Name()), file.Name())
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
}
