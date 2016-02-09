package datastore

import (
	"ape/models"
	"errors"
	"log"
	"os"
	"path/filepath"
)

// ErrNotFound is a not found error
var ErrNotFound = errors.New("git datastore: PkgInfo not Found")

// GitRepo is a munki repo
type GitRepo struct {
	Path          string
	indexPkgsInfo map[string]*models.PkgsInfo
}

// AllPkgsInfos returns a list of pkgsinfos
func (r *GitRepo) AllPkgsInfos() ([]*models.PkgsInfo, error) {
	// load all pkgsinfo from repo
	pkgsinfos, err := loadAllPkgsInfos(r.Path)
	if err != nil {
		return nil, err
	}
	// update index
	r.updatePkgsInfoIndex(pkgsinfos)

	// create an array of *models.PkgsInfo
	var pkgsInfoList []*models.PkgsInfo
	for _, v := range pkgsinfos {
		pkgsInfoList = append(pkgsInfoList, v.PkgsInfo)
	}

	return pkgsInfoList, nil
}

// PkgInfo a single PkgsInfo by name
func (r *GitRepo) PkgInfo(name string) (*models.PkgsInfo, error) {
	// initialize
	pkgsinfos, err := loadAllPkgsInfos(r.Path)
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

func (r *GitRepo) updatePkgsInfoIndex(pkgsinfos []*pkgsInfo) {
	r.indexPkgsInfo = make(map[string]*models.PkgsInfo, len(pkgsinfos))
	for _, info := range pkgsinfos {
		r.indexPkgsInfo[info.Filename] = info.PkgsInfo
	}
}

// PkgsInfo struct
type pkgsInfo struct {
	*models.PkgsInfo
}

// load all pkgsinfos from repo
func loadAllPkgsInfos(path string) ([]*pkgsInfo, error) {
	var pkgsInfoList = new([]*pkgsInfo)
	err := filepath.Walk(path, repoWalkFn(pkgsInfoList))
	if err != nil {
		return nil, err
	}
	return *pkgsInfoList, nil
}

func repoWalkFn(pkgs *[]*pkgsInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		if !info.IsDir() {
			if err := pkginfo(pkgs, info, file); err != nil {
				return err
			}
		}
		return nil
	}
}

// pkginfo adds a new PkgsInfo to the array, but skips if it fails to the code the plist
var pkginfo = func(pkgs *[]*pkgsInfo, info os.FileInfo, file *os.File) error {
	var pkginfo models.PkgsInfo
	err := pkginfo.Decode(file)
	if err != nil {
		log.Printf("git-repo: failed to decode %v, skipping \n", info.Name())
		return nil
	}
	pkginfo.Filename = info.Name()
	*pkgs = append(*pkgs, &pkgsInfo{
		&pkginfo,
	})
	return nil
}
