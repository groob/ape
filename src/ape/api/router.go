package api

import (
	"ape/datastore"
	"fmt"
	"net/http"

	"github.com/groob/plist"
	"github.com/julienschmidt/httprouter"
)

func router(conf *config) http.Handler {
	router := httprouter.New()
	router.GET("/api/manifests", handleManifestsList(conf.db))
	router.GET("/api/manifests/:name", handleManifestsShow(conf.db))
	router.POST("/api/manifests", handleManifestsCreate(conf.db))
	router.DELETE("/api/manifests/:name", handleManifestsDelete(conf.db))
	router.PATCH("/api/manifests/:name", handleManifestsUpdate(conf.db))
	router.GET("/api/pkgsinfos", handlePackagesList(conf.db))
	router.GET("/api/pkgsinfos/:name", handlePackagesShow(conf.db))
	router.POST("/api/pkgsinfos", handlePackagesCreate(conf.db))
	router.DELETE("/api/pkgsinfos/:name", handlePackagesDelete(conf.db))
	router.POST("/api/pkgs/*path", handlePkgsCreate(conf.db))
	router.ServeFiles("/repo/*filepath", http.Dir("repo"))
	router.GET("/catalogs/:name", handleCatalogsShow(conf.db))
	return router
}

func handleCatalogsShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// name := ps.ByName("name")
		pkgsinfos, err := db.AllPkgsInfos()
		if err != nil {
			respondError(rw, http.StatusInternalServerError,
				fmt.Errorf("Failed to fetch pkgsinfo from the datastore: %v", err))
			return
		}

		plist.NewEncoder(rw).Encode(pkgsinfos)
	}
}
