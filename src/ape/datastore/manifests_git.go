package datastore

import (
	"ape/models"
	"fmt"
	"path/filepath"
)

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
