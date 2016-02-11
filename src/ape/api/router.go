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
	router.GET("/api/manifests/:name", handleManifestsShow(conf.db))
	router.POST("/api/manifests", handleManifestsCreate(conf.db))
	router.DELETE("/api/manifests/:name", handleManifestsDelete(conf.db))
	router.PATCH("/api/manifests/:name", handleManifestsUpdate(conf.db))
	return router
}
