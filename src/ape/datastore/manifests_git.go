package datastore

import (
	"ape/models"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var ErrExists = errors.New("Already Exists")

type manifests []models.Manifest

func (m *manifests) add(into decoder) {
	if v, ok := into.(*models.Manifest); ok {
		*m = append(*m, *v)
	}
}

func (m *manifests) load(path string) error {
	manifestPath := fmt.Sprintf("%v/manifests", path)
	err := filepath.Walk(manifestPath, repoWalkFn(m))
	if err != nil {
		return err
	}
	return nil
}

// AllManifests returns an array of manifests
func (r *GitRepo) AllManifests() ([]*models.Manifest, error) {
	m := &manifests{}
	err := m.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updateIndex(m)

	// create an array of manifests
	var manifestList []*models.Manifest
	for _, manifest := range *m {
		v := models.Manifest(manifest)
		manifestList = append(manifestList, &v)
	}
	return manifestList, nil
}

// Manifest returns a single manifest from repo
func (r *GitRepo) Manifest(name string) (*models.Manifest, error) {
	m := &manifests{}
	err := m.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updateIndex(m)
	if _, ok := r.indexManifests[name]; !ok {
		return nil, nil
	}
	return r.indexManifests[name], nil
}

// NewManifest returns a single manifest from repo
func (r *GitRepo) NewManifest(name string) (*models.Manifest, error) {
	manifest := &models.Manifest{
		Filename: name,
	}
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, manifest.Filename)
	// check if exists
	if _, err := os.Stat(manifestPath); err == nil {
		return nil, ErrExists
	}
	// create new
	f, err := os.Create(manifestPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()
	return manifest, nil
}

// SaveManifest saves a manifest to the datastore
func (r *GitRepo) SaveManifest(manifest *models.Manifest) error {
	manifestPath := fmt.Sprintf("%v/manifests/%v", r.Path, manifest.Filename)
	file, err := os.OpenFile(manifestPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := manifest.Encode(file); err != nil {
		return err
	}
	return nil
}
