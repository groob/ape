package datastore

import (
	"ape/models"
	"fmt"
	"path/filepath"
)

// AllPkgsInfos returns a list of pkgsinfos
func (r *GitRepo) AllPkgsInfos() ([]*models.PkgsInfo, error) {
	// load all pkgsinfo from repo
	pkgsinfos := pkgsInfos{}
	err := pkgsinfos.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updatePkgsInfoIndex(pkgsinfos)

	// create an array of *models.PkgsInfo
	var pkgsInfoList []*models.PkgsInfo
	for _, info := range pkgsinfos {
		v := models.PkgsInfo(info)
		pkgsInfoList = append(pkgsInfoList, &v)
	}
	return pkgsInfoList, nil
}

// PkgInfo a single PkgsInfo by name
func (r *GitRepo) PkgInfo(name string) (*models.PkgsInfo, error) {
	// initialize
	pkgsinfos := pkgsInfos{}
	err := pkgsinfos.load(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updatePkgsInfoIndex(pkgsinfos)

	if _, ok := r.indexPkgsInfo[name]; !ok {
		return nil, nil
	}
	return r.indexPkgsInfo[name], nil
}

func (r *GitRepo) updatePkgsInfoIndex(pkgsinfos pkgsInfos) {
	r.indexPkgsInfo = make(map[string]*models.PkgsInfo, len(pkgsinfos))
	for _, info := range pkgsinfos {
		v := models.PkgsInfo(info)
		r.indexPkgsInfo[info.Filename] = &v
	}
}

type pkgsInfos []models.PkgsInfo

func (p *pkgsInfos) add(into decoder) {
	if v, ok := into.(*models.PkgsInfo); ok {
		*p = append(*p, *v)
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
