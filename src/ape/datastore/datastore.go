package datastore

import "ape/models"

// Datastore is an interface around munki storage
type Datastore interface {
	AllPkgsInfos() ([]*models.PkgsInfo, error)
	PkgInfo(name string) (*models.PkgsInfo, error)
}
