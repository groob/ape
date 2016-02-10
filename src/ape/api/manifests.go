package api

import (
	"ape/datastore"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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
		f := &struct {
			Name string `json:"name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(f)
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusNoContent)
			return
		default:
			log.Println(err)
		}

		manifest, err := db.NewManifest(f.Name)
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
