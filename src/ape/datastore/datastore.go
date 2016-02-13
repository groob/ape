package datastore

import (
	"ape/models"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ErrExists file already exists
var ErrExists = errors.New("Resource already exists")

// ErrNotFound = resource not found
var ErrNotFound = errors.New("Resource not found")

// Datastore is an interface around munki storage
type Datastore interface {
	manifestStore
	pkgsinfoStore
	pkgsStore
}

type manifestStore interface {
	AllManifests() (*models.ManifestCollection, error)
	Manifest(name string) (*models.Manifest, error)
	NewManifest(name string) (*models.Manifest, error)
	SaveManifest(manifest *models.Manifest) error
	DeleteManifest(name string) error
}

type pkgsinfoStore interface {
	AllPkgsinfos() (*models.PkgsInfoCollection, error)
	Pkgsinfo(name string) (*models.PkgsInfo, error)
	NewPkgsinfo(name string) (*models.PkgsInfo, error)
	SavePkgsinfo(manifest *models.PkgsInfo) error
	DeletePkgsinfo(name string) error
}

type pkgsStore interface {
	NewPkg(filename string, body io.Reader) error
	DeletePkg(name string) error
}

// SimpleRepo is a filesystem based backend
type SimpleRepo struct {
	Path           string
	indexManifests map[string]*models.Manifest
	indexPkgsinfo  map[string]*models.PkgsInfo
}

// NewPkg creates a new pkg file
func (r *SimpleRepo) NewPkg(filename string, body io.Reader) error {
	pkgPath := fmt.Sprintf("%v/pkgs/%v", r.Path, filename)
	// check if exists
	if _, err := os.Stat(pkgPath); err == nil {
		return ErrExists
	}

	dir := filepath.Dir(pkgPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// create new
	f, err := os.Create(pkgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, body)
	if err != nil {
		return err
	}
	return nil
}

// DeletePkg deletes a pkg file
func (r *SimpleRepo) DeletePkg(name string) error {
	pkgPath := fmt.Sprintf("%v/pkgs/%v", r.Path, name)
	return os.Remove(pkgPath)
}
