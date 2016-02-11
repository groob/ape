package models

import "errors"

// ErrNotFound when a resource is empy
var ErrNotFound = errors.New("Not Found")

// Viewer interface
type Viewer interface {
	View() (*APIResponse, error)
}

// APIResponse ...
type APIResponse struct {
	Data []byte
}
