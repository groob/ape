package datastore

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/groob/ape/models"
)

var (
	// ErrExists file already exists
	ErrExists = errors.New("Resource already exists")

	// ErrNotFound = resource not found
	ErrNotFound = errors.New("Resource not found")
)

// Datastore is an interface around munki storage
type Datastore interface {
	models.PkgsinfoStore
	models.PkgStore
	models.ManifestStore
}

// SimpleRepo is a filesystem based backend
type SimpleRepo struct {
	Path           string
	indexManifests map[string]*models.Manifest
	indexPkgsinfo  map[string]*models.PkgsInfo
}

func deleteFile(path string) error {
	if err := os.Remove(path); err != nil {
		return ErrNotFound
	}
	return nil
}

func createFile(path string) error {
	// check if exists
	if _, err := os.Stat(path); err == nil {
		return ErrExists
	}
	// create the relative directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
