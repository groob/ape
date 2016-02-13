package api

import (
	"ape/datastore"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func handlePkgsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		path := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.NewPkg(path, r.Body)
		if _, ok := err.(*os.PathError); ok {
			respondError(rw, http.StatusNotFound, accept, err)
			return
		}
		if err != nil {
			respondError(rw, http.StatusInternalServerError, accept, err)
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("done"))
		return
	}
}

func handlePkgsDelete(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := strings.TrimLeft(ps.ByName("name"), "/")
		accept := acceptHeader(r)
		err := db.DeletePkg(name)
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
