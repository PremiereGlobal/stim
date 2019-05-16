package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
)

// CreateFileIfNotExist returns whether the given file exists
// Attempts to create the full path and file if it does not exist
func CreateFileIfNotExist(filePath string) error {
	log := stimlog.GetLogger()

	// Frist check to see if given filePath isn't a directory
	isDir, _ := IsDirectory(filePath)
	if isDir == true {
		return errors.New("given file path is a directory and not a path to a file")
	}

	// Check and create the base path if needed
	dir, _ := filepath.Split(filePath)
	if len(dir) > 0 {
		err := CreateDirIfNotExist(dir)
		if err != nil {
			return err
		}
	}

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Debug("Creating file: '" + filePath + "'")
		f, err := os.Create(filePath)
		defer f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateDirIfNotExist returns whether the given directory exists
// Attempts to create the full path if it does not exist
func CreateDirIfNotExist(path string) error {
	log := stimlog.GetLogger()

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		log.Debug("Creating folder: '" + path + "'")
		err := os.MkdirAll(path, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return nil
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
