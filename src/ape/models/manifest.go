package models

import (
	"io"

	"github.com/groob/plist"
)

// Manifest represents the structure of a munki manifest plist
type Manifest struct {
	Filename          string   `plist:"-" json:"-"`
	Catalogs          []string `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DisplayName       string   `plist:"display_name,omitempty" json:"display_name,omitempty"`
	IncludedManifests []string `plist:"included_manifests,omitempty" json:"included_manifests,omitempty"`
	ManagedInstalls   []string `plist:"managed_installs,omitempty" json:"managed_installs,omitempty"`
}

// Decode a plist into a struct
func (m *Manifest) Decode(r io.Reader) error {
	return plist.NewDecoder(r).Decode(m)
}

// Encode a go struct into a plist
func (m *Manifest) Encode(w io.Writer) error {
	enc := plist.NewEncoder(w)
	enc.Indent("  ")
	return enc.Encode(m)
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

// ManifestListx returns a default view in api response
func ManifestListx(manifests []*Manifest) *View {
	var viewDefault = make(View, len(manifests))
	for _, info := range manifests {
		viewDefault[info.Filename] = info

	}
	return &viewDefault
}

// ManifestList represents a slice of manifests
type ManifestList []*Manifest

// View returns a view
func (l ManifestList) View() *View {
	var viewDefault = make(View, len(l))
	for _, info := range l {
		viewDefault[info.Filename] = info

	}
	return &viewDefault
}
