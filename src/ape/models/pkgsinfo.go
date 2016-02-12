package models

import (
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/groob/plist"
)

var errNoContent = errors.New("No Content")

// PkgsInfo represents the structure of a pkgsinfo file
type PkgsInfo struct {
	Filename              string        `plist:"-" json:"name"`
	Metadata              metadata      `plist:"_metadata" json:"_metadata"`
	Autoremove            bool          `plist:"autoremove" json:"autoremove"`
	Catalogs              []string      `plist:"catalogs" json:"catalogs"`
	Description           string        `plist:"description" json:"description"`
	DisplayName           string        `plist:"display_name" json:"display_name"`
	InstallerItemHash     string        `plist:"installer_item_hash" json:"installer_item_hash"`
	InstallerItemLocation string        `plist:"installer_item_location" json:"installer_item_location"`
	InstallerItemSize     int           `plist:"installer_item_size" json:"installer_item_size"`
	InstallerType         string        `plist:"installer_type" json:"installer_item_type"`
	Installs              []installs    `plist:"installs" json:"installs"`
	ItemsToCopy           []itemsToCopy `plist:"items_to_copy" json:"items_to_copy"`
	MinimumOsVersion      string        `plist:"minimum_os_version" json:"minimum_os_version"`
	Name                  string        `plist:"name" json:"name"`
	UnattendedInstall     bool          `plist:"unattended_install" json:"unattended_install"`
	AppleItem             bool          `plist:"apple_item" json:"apple_item"`
	BlockingApplications  []string      `plist:"blocking_applications" json:"blocking_applications"`
}

type metadata struct {
	CreatedBy    string    `plist:"created_by" json:"created_by"`
	CreatedDate  time.Time `plist:"creation_date" json:"created_date"`
	MunkiVersion string    `plist:"munki_version" json:"munki_version"`
	OSVersion    string    `plist:"os_version" json:"os_version"`
}

type installs struct {
	CFBundleIdentifier         string `plist:"CFBundleIdentifier"`
	CFBundleName               string `plist:"CFBundleName"`
	CFBundleShortVersionString string `plist:"CFBundleShortVersionString"`
	CFBundleVersion            string `plist:"CFBundleVersion"`
	MinOSVersion               string `plist:"minosversion" json:"min_os_version"`
	Path                       string `plist:"path" json:"path"`
	Type                       string `plist:"type" json:"type"`
	VersionComparisonKey       string `plist:"version_comparison_key" json:"version_comparision_key"`
}

type itemsToCopy struct {
	DestinationPath string `plist:"destination_path" json:"destination_path"`
	SourceItem      string `plist:"source_item" json:"source_item"`
}

// Decode a plist into a struct
func (p *PkgsInfo) Decode(r io.Reader) error {
	return plist.NewDecoder(r).Decode(p)
}

// Encode a go struct into a plist
func (p *PkgsInfo) Encode(w io.Writer) error {
	enc := plist.NewEncoder(w)
	enc.Indent("  ")
	return enc.Encode(p)
}

// View returns an pkgsinfo api response
func (p *PkgsInfo) View() (*APIResponse, error) {
	if p == nil {
		return nil, ErrNotFound
	}
	response := &APIResponse{}
	data, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return response, err
	}
	response.Data = data
	return response, nil
}

type defaultPkgsInfoView struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Catalogs    []string `json:"catalogs"`
	PkgURL      string   `json:"pkg_url,omitempty"`
}

// PkgsInfoList is a list of pkgsinfos
type PkgsInfoList []*PkgsInfo

// View returns an view
func (p *PkgsInfoList) View() (*APIResponse, error) {
	var listView []defaultPkgsInfoView
	for _, info := range *p {
		listView = append(listView, defaultPkgsInfoView{
			Name:        info.Filename,
			DisplayName: info.DisplayName,
			Catalogs:    info.Catalogs,
			PkgURL:      info.InstallerItemLocation,
		})
	}
	response := &APIResponse{}
	data, err := json.MarshalIndent(listView, "", " ")
	if err != nil {
		return response, err
	}
	response.Data = data
	return response, nil
}
