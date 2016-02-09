package models

import (
	"io"

	"github.com/groob/plist"
)

// Manifest represents the structure of a munki manifest plist
type Manifest struct {
	Filename          string   `plist:"-" json:"-"`
	Catalogs          []string `plist:"catalogs" json:"catalogs"`
	DisplayName       string   `plist:"display_name" json:"display_name"`
	IncludedManifests []string `plist:"included_manifests" json:"included_manifests"`
	ManagedInstalls   []string `plist:"managed_installs" json:"managed_installs"`
}

// Decode a plist into a struct
func (m *Manifest) Decode(r io.Reader) error {
	return plist.NewDecoder(r).Decode(m)
}

// Encode a go struct into a plist
func (m *Manifest) Encode(w io.Writer) error {
	return plist.NewEncoder(w).Encode(m)
}

// View returns a map for the JSON response
func (m *Manifest) View() *View {
	if m == nil {
		return &View{}
	}
	var view = make(View, 0)
	view[m.Filename] = m
	return &view
}

// ManifestList returns a default view in api response
func ManifestList(manifests []*Manifest) *View {
	var viewDefault = make(View, len(manifests))
	for _, info := range manifests {
		viewDefault[info.Filename] = info

	}
	return &viewDefault
}
