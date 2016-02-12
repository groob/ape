package datastore

import (
	"ape/models"
	"fmt"
	"log"
	"os"
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

// NewPkgInfo returns a single manifest from repo
func (r *GitRepo) NewPkgInfo(name string) (*models.PkgsInfo, error) {
	pkgsinfo := &models.PkgsInfo{
		Filename: name,
	}
	pkgsinfosPath := fmt.Sprintf("%v/pkgsinfos/%v", r.Path, pkgsinfo.Filename)
	// check if exists
	if _, err := os.Stat(pkgsinfosPath); err == nil {
		return nil, ErrExists
	}
	// create new
	f, err := os.Create(pkgsinfosPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()
	return pkgsinfo, nil
}

// SavePkgInfo saves a manifest to the datastore
func (r *GitRepo) SavePkgInfo(pkgsinfo *models.PkgsInfo) error {
	pkgsinfosPath := fmt.Sprintf("%v/pkgsinfos/%v", r.Path, pkgsinfo.Filename)
	file, err := os.OpenFile(pkgsinfosPath, os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := pkgsinfo.Encode(file); err != nil {
		return err
	}
	return nil
}

// DeletePkgInfo ...
func (r *GitRepo) DeletePkgInfo(name string) error {
	pkgsinfoPath := fmt.Sprintf("%v/pksinfos/%v", r.Path, name)
	err := os.Remove(pkgsinfoPath)
	if err != nil {
		return err
	}
	return nil
}

type pkgsInfos models.PkgsInfoList

// implement loader interface
func (p *pkgsInfos) add(into decoder) {
	if v, ok := into.(*models.PkgsInfo); ok {
		*p = append(*p, v)
	}
}

// implement loader interface
func (p *pkgsInfos) load(path string) error {
	pkgsInfosPath := fmt.Sprintf("%v/pkgsinfo", path)
	err := filepath.Walk(pkgsInfosPath, repoWalkFn(p))
	if err != nil {
		return err
	}
	return nil
}
