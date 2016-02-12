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

	"github.com/julienschmidt/httprouter"
)

func handlePackagesList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		pkgsinfos, err := db.AllPkgsInfos()
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch PkgsInfo list from the datastore: %v", err))
			return
		}
		respondOK(rw, pkgsinfos)
	}
}

//
func handlePackagesShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		pkgsinfo, err := db.PkgInfo(name)
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch pkgsinfo from the datastore: %v", err))
			return
		}
		respondOK(rw, pkgsinfo)
	}
}

func handlePackagesCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		payload := &models.PkgsInfo{}
		err := json.NewDecoder(r.Body).Decode(payload)
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusBadRequest)
			return
		default:
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to create new pkgsinfo: %v", err))
			return
		}
		var pkgsinfo *models.PkgsInfo
		if payload.Filename == "" {
			// filename is required
			respondError(rw, http.StatusBadRequest,
				errors.New("the name field is required to create a pkgsinfo"))
			return
		}

		// If the body contains a valid pkgsinfo, create it
		pkgsinfo = payload

		_, err = db.NewPkgInfo(payload.Filename)
		switch err {
		case nil:
			break
		case datastore.ErrExists:
			respondError(rw, http.StatusConflict, err)
			return
		default:
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to create new pkgsinfo: %v", err))
			return
		}

		if err := db.SavePkgInfo(pkgsinfo); err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to save pkgsinfo: %v", err))
			return
		}
		respondCreated(rw, pkgsinfo, "")
	}
}

func handlePackagesDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		err := db.DeletePkgInfo(name)

		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to delete pkgsinfo from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
