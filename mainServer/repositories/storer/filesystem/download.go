package filesystem

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

const IllegalChars = "<>:\"/\\|?*"

// MakeDownloadZip copies the all contents of the source directory to a .zip file with the specified
func (fs FileSystem) MakeDownloadZip(filename string, sourceDir string) (string, error) {
	path, err := fs.GetDownloadPath(filename)
	if err != nil {
		return "", err
	}

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

	// write directory contents to the writer
	err = addFilesInDirToZip(zipWriter, sourceDir, "")
	if err != nil {
		return "", err
	}

	return path, nil
}

// write files in dirPath to the location dirInZip in the zipWriter
func addFilesInDirToZip(zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Wrapped in function to allow for "defer file.close()"
		err := addFileInDirToZip(file, zipWriter, dirPath, dirInZip)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFileInDirToZip(file fs.FileInfo, zipWriter *zip.Writer, dirPath string, dirInZip string) error {
	if file.IsDir() {
		//Check if it is not the git folder
		if file.Name() != ".git" {
			err := addFilesInDirToZip(zipWriter, filepath.Join(dirPath, file.Name()), file.Name())
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
