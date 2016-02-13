package models

type catalogInfo struct {
	Catalogs    []string `plist:"catalogs" json:"catalogs"`
	Category    string   `plist:"category,omitempty"`
	Description string   `plist:"description" json:"description"`
	Name        string   `plist:"name" json:"name"`
}

// Catalog represents a munki catalog
type Catalog []*catalogInfo
