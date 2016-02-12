package api

import (
	"ape/datastore"
	"fmt"
	"net/http"

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
