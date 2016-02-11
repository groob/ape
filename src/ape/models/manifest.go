package models

import (
	"encoding/json"
	"io"

	"github.com/groob/plist"
)

// Manifest represents the structure of a munki manifest plist
type Manifest struct {
	Filename          string   `plist:"-" json:"name"`
	Catalogs          []string `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DisplayName       string   `plist:"display_name,omitempty" json:"display_name,omitempty"`
	IncludedManifests []string `plist:"included_manifests,omitempty" json:"included_manifests,omitempty"`
	OptionalInstalls  []string `plist:"optional_installs,omitempty" json:"optional_installs,omitempty"`
	ManagedInstalls   []string `plist:"managed_installs,omitempty" json:"managed_installs,omitempty"`
	Notes             string   `plist:"notes,omitempty" json:"notes,omitempty"`
	User              string   `plist:"user,omitempty" json:"user,omitempty"`
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
func (m *Manifest) View() (*APIResponse, error) {
	if m == nil {
		return nil, ErrNotFound
	}
	response := &APIResponse{}
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return response, err
	}

	response.Data = data
	return response, nil
}

// ManifestList represents a slice of manifests
type ManifestList []*Manifest

// View returns a view
func (m *ManifestList) View() (*APIResponse, error) {
	response := &APIResponse{}
	data, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return response, err
	}
	response.Data = data
	return response, nil
}
