package datastore

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// NewPkg creates a new pkg file
func (r *GitRepo) NewPkg(filename string, body io.Reader) error {
	pkgPath := fmt.Sprintf("%v/pkgs/%v", r.Path, filename)
	// check if exists
	if _, err := os.Stat(pkgPath); err == nil {
		return ErrExists
	}

	dir := filepath.Dir(pkgPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// create new
	f, err := os.Create(pkgPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, body)
	if err != nil {
		return err
	}
	return nil
}
