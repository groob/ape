package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func router(conf *config) http.Handler {
	router := httprouter.New()
	router.GET("/api/packages", handlePackagesList(conf.db))
	router.GET("/api/packages/:name", handlePackagesShow(conf.db))
	router.GET("/api/manifests", handleManifestsList(conf.db))
	return router
}
