package datastore

import (
	"ape/models"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

// ErrExists file already exists
var ErrExists = errors.New("Resource already exists")

// Datastore is an interface around munki storage
type Datastore interface {
	pkgsinfoStore
	manifestStore
}

type pkgsinfoStore interface {
	AllPkgsInfos() (*models.PkgsInfoList, error)
	PkgInfo(name string) (*models.PkgsInfo, error)
	NewPkgInfo(name string) (*models.PkgsInfo, error)
	SavePkgInfo(pkginfo *models.PkgsInfo) error
	DeletePkgInfo(name string) error
}

type manifestStore interface {
	AllManifests() (*models.ManifestList, error)
	Manifest(name string) (*models.Manifest, error)
	NewManifest(name string) (*models.Manifest, error)
	SaveManifest(manifest *models.Manifest) error
	DeleteManifest(name string) error
}

// GitRepo is a munki repo
type GitRepo struct {
	Path           string
	indexPkgsInfo  map[string]*models.PkgsInfo
	indexManifests map[string]*models.Manifest
}

func (r *GitRepo) updateIndex(l loader) {
	switch l {
	case l.(*pkgsInfos):
		pkgsinfos := l.(*pkgsInfos)
		r.indexPkgsInfo = make(map[string]*models.PkgsInfo, len(*pkgsinfos))
		for _, info := range *pkgsinfos {
			v := models.PkgsInfo(*info)
			r.indexPkgsInfo[info.Filename] = &v
		}
	case l.(*manifests):
		manifests := l.(*manifests)
		r.indexManifests = make(map[string]*models.Manifest, len(*manifests))
		for _, info := range *manifests {
			v := models.Manifest(*info)
			r.indexManifests[info.Filename] = &v
		}
	}
}

type loader interface {
	load(string) error
	add(decoder)
}

type decoder interface {
	Decode(io.Reader) error
}

func repoWalkFn(pkgs loader) filepath.WalkFunc {
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
			if err := open(pkgs, info, file); err != nil {
				return err
			}
		}
		return nil
	}
}

// pkginfo adds a new PkgsInfo to the array, but skips if it fails to the code the plist
var open = func(l loader, info os.FileInfo, file *os.File) error {
	var into decoder
	if _, ok := l.(*pkgsInfos); ok {
		into = &models.PkgsInfo{}
	}
	if _, ok := l.(*manifests); ok {
		into = &models.Manifest{}
	}

	err := into.Decode(file)
	if err != nil {
		log.Printf("git-repo: failed to decode %v, skipping \n", info.Name())
		return nil
	}
	switch into {
	case into.(*models.PkgsInfo):
		v := into.(*models.PkgsInfo)
		v.Filename = info.Name()
		l.add(v)
		return nil
	case into.(*models.Manifest):
		v := into.(*models.Manifest)
		v.Filename = info.Name()
		l.add(v)
		return nil
	default:
		return errors.New("git datastore: wrong decoder type")
	}
}
