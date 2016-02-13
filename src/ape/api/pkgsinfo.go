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

func handlePkgsinfoList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		pkgsinfos, err := db.AllPkgsinfos()
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to fetch pkgsinfo list from the datastore: %v", err))
			return
		}
		// apply any filters
		pkgsinfos = applyPkgsinfoFilters(pkgsinfos, r.URL.Query())
		respondOK(rw, pkgsinfos, accept)
	}
}

func handlePkgsinfoShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		pkgsinfo, err := db.Pkgsinfo(name)
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to fetch pkgsinfo from the datastore: %v", err))
			return
		}
		respondOK(rw, pkgsinfo, accept)
	}
}

func handlePkgsinfoCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		accept := acceptHeader(r)
		contentType := contentHeader(r)
		var payload struct {
			Filename string `plist:"filename" json:"filename"`
			*models.PkgsInfo
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
			rw.WriteHeader(http.StatusBadRequest)
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to create new pkgsinfo: %v", err))
			return
		}
		var pkgsinfo *models.PkgsInfo

		if payload.Filename == "" {
			// filename is required
			respondError(rw, http.StatusBadRequest, accept,
				errors.New("the name field is required to create a pkgsinfo"))
			return
		}

		// If the body contains a valid pkgsinfo, create it
		pkgsinfo = payload.PkgsInfo
		pkgsinfo.Filename = payload.Filename

		_, err = db.NewPkgsinfo(payload.Filename)
		switch err {
		case nil:
			break
		case datastore.ErrExists:
			respondError(rw, http.StatusConflict, accept, err)
			return
		default:
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to create new pkgsinfo: %v", err))
			return
		}

		if err := db.SavePkgsinfo(pkgsinfo); err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to save pkgsinfo: %v", err))
			return
		}
		respondCreated(rw, pkgsinfo, accept)
	}
}

func handlePkgsinfoDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeletePkgsinfo(name)
		// path error, return not found
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}

		// all other errors
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to delete pkgsinfo from the datastore: %v", err))
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func handlePkgsinfoReplace(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		contentType := contentHeader(r)
		pkgsinfo, err := db.Pkgsinfo(name)
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
				fmt.Errorf("Failed to update pkgsinfo from the datastore: %v", err))
			return
		}

		payload := &models.PkgsInfo{}
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
				fmt.Errorf("Failed to update pkgsinfo: %v", err))
			return
		}

		if err := db.SavePkgsinfo(pkgsinfo); err != nil {
			respondError(rw, http.StatusInternalServerError, accept,
				fmt.Errorf("Failed to save pkginfo: %v", err))
			return
		}

		pkgsinfo = payload
		pkgsinfo.Filename = name

		// manifest updated ok, respond
		respondOK(rw, pkgsinfo, accept)
	}
}
