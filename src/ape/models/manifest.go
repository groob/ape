package models

// ManifestStore is the interface for accessing manifests in a database or filesystem
type ManifestStore interface {
	AllManifests() (*ManifestCollection, error)
	Manifest(name string) (*Manifest, error)
	NewManifest(name string) (*Manifest, error)
	SaveManifest(manifest *Manifest) error
	DeleteManifest(name string) error
}

// Manifest represents the structure of a munki manifest
// This is what would be serialized in a datastore
type Manifest struct {
	Filename          string      `plist:"-" json:"-"`
	Catalogs          []string    `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DisplayName       string      `plist:"display_name,omitempty" json:"display_name,omitempty"`
	IncludedManifests []string    `plist:"included_manifests,omitempty" json:"included_manifests,omitempty"`
	Notes             string      `plist:"notes,omitempty" json:"notes,omitempty"`
	User              string      `plist:"user,omitempty" json:"user,omitempty"`
	ConditionalItems  []condition `plist:"conditional_items,omitempty" json:"conditional_items,omitempty"`
	manifestItems
}

type manifestItems struct {
	OptionalInstalls  []string `plist:"optional_installs,omitempty" json:"optional_installs,omitempty"`
	ManagedInstalls   []string `plist:"managed_installs,omitempty" json:"managed_installs,omitempty"`
	ManagedUninstalls []string `plist:"managed_uninstalls,omitempty" json:"managed_uninstalls,omitempty"`
	ManagedUpdates    []string `plist:"managed_updates,omitempty" json:"managed_updates,omitempty"`
}

type condition struct {
	Condition string `plist:"condition" json:"condition"`
	manifestItems
}

// ManifestView is the response view
type manifestView struct {
	Filename string `plist:"filename,omitempty" json:"filename,omitempty"`
	*Manifest
}

// View returns response
func (m *Manifest) View(accept string) (*Response, error) {
	if m == nil {
		return nil, ErrNoData
	}

	return marshal(m, accept)
}

// ManifestCollection represents a list of manifests
type ManifestCollection []*Manifest

// View returns response
func (m *ManifestCollection) View(accept string) (*Response, error) {
	var view []*manifestView
	for _, item := range *m {
		viewItem := &manifestView{
			item.Filename,
			item,
		}
		view = append(view, viewItem)
	}

	return marshal(view, accept)
}
