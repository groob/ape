package models

import "encoding/json"

// Viewer interface
type Viewer interface {
	View() *View
}

// View model
type View map[string]interface{}

// ToJSON returns json array or error
func (v *View) ToJSON() ([]byte, error) {
	if len(*v) == 0 {
		return nil, errNoContent
	}
	return json.MarshalIndent(v, "", " ")
}
