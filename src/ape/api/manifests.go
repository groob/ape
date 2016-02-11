package api

import (
	"ape/datastore"
	"ape/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func handleManifestsList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		manifests, err := db.AllManifests()
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch manifest list from the datastore: %v", err))
			return
		}
		respondOK(rw, manifests)
	}
}

func handleManifestsShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		manifest, err := db.Manifest(name)
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch manifest from the datastore: %v", err))
			return
		}
		respondOK(rw, manifest)
	}
}

func handleManifestsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		payload := &models.Manifest{}
		err := json.NewDecoder(r.Body).Decode(payload)
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusBadRequest)
			return
		default:
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		}
		var manifest *models.Manifest

		if payload.Filename == "" {
			// filename is required
			respondError(rw, http.StatusBadRequest,
				errors.New("the name field is required to create a manifest"))
			return
		}

		// If the body contains a valid manifest, create it
		manifest = payload

		_, err = db.NewManifest(payload.Filename)
		switch err {
		case nil:
			break
		case datastore.ErrExists:
			respondError(rw, http.StatusConflict, err)
			return
		default:
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to create new manifest: %v", err))
			return
		}

		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}
		respondCreated(rw, manifest, "")
	}
}

func handleManifestsDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		err := db.DeleteManifest(name)

		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to delete manifest from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func handleManifestsUpdate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		manifest, err := db.Manifest(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch manifest from the datastore: %v", err))
			return
		}

		// succesfully retrieved manifest, deal with payload
		payload := &manifestPayload{}
		err = json.NewDecoder(r.Body).Decode(payload)
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusBadRequest)
			return
		default:
			log.Println(err)
		}

		// filename must match
		if *payload.Filename != name {
			respondError(rw, http.StatusBadRequest, nil)
			return
		}

		if payload.Catalogs != nil {
			manifest.Catalogs = *payload.Catalogs
		}

		if payload.DisplayName != nil {
			manifest.DisplayName = *payload.DisplayName
		}

		// no save
		if err := db.SaveManifest(manifest); err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to save manifest: %v", err))
			return
		}

		// manifest updated ok, respond
		respondOK(rw, manifest)
	}
}

// use manifestPayload to check for nil values during an update
type manifestPayload struct {
	Filename          *string   `plist:"-" json:"name"`
	Catalogs          *[]string `plist:"catalogs,omitempty" json:"catalogs,omitempty"`
	DisplayName       *string   `plist:"display_name,omitempty" json:"display_name,omitempty"`
	IncludedManifests *[]string `plist:"included_manifests,omitempty" json:"included_manifests,omitempty"`
	OptionalInstalls  *[]string `plist:"optional_installs,omitempty" json:"optional_installs,omitempty"`
	ManagedInstalls   *[]string `plist:"managed_installs,omitempty" json:"managed_installs,omitempty"`
	Notes             *string   `plist:"notes,omitempty" json:"notes,omitempty"`
	User              *string   `plist:"user,omitempty" json:"user,omitempty"`
}
