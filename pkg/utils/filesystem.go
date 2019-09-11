package utils

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
)

const (
	UserOnlyMode  os.FileMode = 0700
	UserGroupMode os.FileMode = 0700
	WorldMode     os.FileMode = 0700
)

// CreateFileIfNotExist returns whether the given file exists
// Attempts to create the full path and file if it does not exist
func CreateFileIfNotExist(filePath string, perm os.FileMode) error {
	log := stimlog.GetLogger()

	stat, err := os.Stat(filePath)
	if err == nil && !stat.IsDir() {
		return nil
	} else if err == nil && stat.IsDir() {
		return errors.New("given file path is a directory and not a path to a file")
	}

	// Check and create the base path if needed
	dir, _ := filepath.Split(filePath)
	if len(dir) > 0 {
		err := CreateDirIfNotExist(dir, perm)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		log.Debug("Creating file: '{}", filePath)
		f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// CreateDirIfNotExist returns whether the given directory exists
// Attempts to create the full path if it does not exist
func CreateDirIfNotExist(path string, perm os.FileMode) error {
	log := stimlog.GetLogger()

	s, err := os.Stat(path)
	if err == nil && s.IsDir() {
		return nil
	} else if err == nil && !s.IsDir() {
		return errors.New("Path is already a file!")
	}
	if !os.IsNotExist(err) {
		return err
	}
	log.Debug("Creating folder: '{}'", path)
	err = os.MkdirAll(path, perm)
	if err == nil {
		return nil
	}
	return err
}

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}
