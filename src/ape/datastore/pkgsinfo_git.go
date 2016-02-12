package datastore

import (
	"ape/models"
	"fmt"
	"path/filepath"
)

// AllPkgsInfos returns a list of pkgsinfos
func (r *GitRepo) AllPkgsInfos() (*models.PkgsInfoList, error) {
	p := &pkgsInfos{}
	err := p.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updateIndex(p)

	list := models.PkgsInfoList(*p)

	return &list, nil
}

// PkgInfo a single PkgsInfo by name
func (r *GitRepo) PkgInfo(name string) (*models.PkgsInfo, error) {
	p := &pkgsInfos{}
	err := p.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updateIndex(p)
	pkgsinfo, ok := r.indexPkgsInfo[name]
	if !ok {
		return nil, nil
	}
	return pkgsinfo, nil
}

type pkgsInfos models.PkgsInfoList

func (p *pkgsInfos) add(into decoder) {
	if v, ok := into.(*models.PkgsInfo); ok {
		*p = append(*p, v)
	}
}

func (p *pkgsInfos) load(path string) error {
	pkgsInfosPath := fmt.Sprintf("%v/pkgsinfo", path)
	err := filepath.Walk(pkgsInfosPath, repoWalkFn(p))
	if err != nil {
		return err
	}
	return nil
}
