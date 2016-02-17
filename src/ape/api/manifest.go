package api

import (
	"ape/datastore"
	"ape/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/groob/plist"
	"github.com/julienschmidt/httprouter"
)

func handleManifestsList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		manifests, err := db.AllManifests()
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to fetch manifest list from the datastore: %v", err))
			return
		}
		respondOK(rw, manifests, accept)
	}
}

func handleManifestsShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		manifest, err := db.Manifest(name)
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to fetch manifest from the datastore: %v", err))
			return
		}
		respondOK(rw, manifest, accept)
	}
}

func handleManifestsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		contentType := contentHeader(r)
		var payload struct {
			Filename string `plist:"filename" json:"filename"`
			*models.Manifest
		}
		var err error

		// decode into xml or json
		switch contentType {
		case "application/xml":
			err = plist.NewDecoder(r.Body).Decode(&payload)
		case "application/json":
			err = json.NewDecoder(r.Body).Decode(&payload)
		default:
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// check error
		switch err {
		case nil:
			break
		case io.EOF:
			respondError(rw, http.StatusBadRequest, accept,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		}
		var manifest *models.Manifest

		if payload.Filename == "" {
			// filename is required
			respondError(rw, http.StatusBadRequest, accept,
				errors.New("the name field is required to create a manifest"))
			return
		}

		// If the body contains a valid manifest, create it
		manifest = payload.Manifest
		manifest.Filename = payload.Filename

		_, err = db.NewManifest(payload.Filename)
		switch err {
		case nil:
			break
		case datastore.ErrExists:
			respondError(rw, http.StatusConflict, accept, err)
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		}

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}
		respondCreated(rw, manifest, accept)
	}
}

func handleManifestsDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeleteManifest(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func handleManifestsUpdate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		contentType := contentHeader(r)
		manifest, err := db.Manifest(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}

		// handle not found err
		if err != nil && err == datastore.ErrNotFound {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}
		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}

		payload := &models.ManifestPayload{}

		// decode into xml or json
		switch contentType {
		case "application/xml":
			err = plist.NewDecoder(r.Body).Decode(payload)
		case "application/json":
			err = json.NewDecoder(r.Body).Decode(payload)
		default:
			respondError(rw, http.StatusBadRequest, accept,
				fmt.Errorf("Failed to update manifest: %v", payload))
			return
		}

		// handle err
		switch err {
		case nil:
			break
		case io.EOF:
			respondError(rw, http.StatusBadRequest, accept,
				fmt.Errorf("Failed to update manifest: %v", err))
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to update manifest: %v", err))
			return
		}

		if payload.Catalogs != nil {
			manifest.Catalogs = *payload.Catalogs
		}

		if payload.DisplayName != nil {
			manifest.DisplayName = *payload.DisplayName
		}

		if payload.IncludedManifests != nil {
			manifest.IncludedManifests = *payload.IncludedManifests
		}

		if payload.OptionalInstalls != nil {
			manifest.OptionalInstalls = *payload.OptionalInstalls
		}

		if payload.ManagedInstalls != nil {
			manifest.ManagedInstalls = *payload.ManagedInstalls
		}

		if payload.ManagedUninstalls != nil {
			manifest.ManagedUninstalls = *payload.ManagedUninstalls
		}

		if payload.ManagedUpdates != nil {
			manifest.ManagedUpdates = *payload.ManagedUpdates
		}

		if payload.Notes != nil {
			manifest.Notes = *payload.Notes
		}

		if payload.User != nil {
			manifest.User = *payload.User
		}

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		// manifest updated ok, respond
		respondOK(rw, manifest, accept)
	}
}

func handleManifestsReplace(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		contentType := contentHeader(r)
		manifest, err := db.Manifest(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}

		// handle not found err
		if err != nil && err == datastore.ErrNotFound {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}
		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}

		payload := &models.Manifest{}
		// decode into xml or json
		switch contentType {
		case "application/xml":
			err = plist.NewDecoder(r.Body).Decode(payload)
		case "application/json":
			err = json.NewDecoder(r.Body).Decode(payload)
		default:
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		// handle err
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusBadRequest)
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to update manifest: %v", err))
			return
		}

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		manifest = payload
		manifest.Filename = name

		// manifest updated ok, respond
		respondOK(rw, manifest, accept)
	}
}
