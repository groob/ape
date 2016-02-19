package datastore

import (
	"ape/models"
	"errors"
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
