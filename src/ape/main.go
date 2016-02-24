package main

import (
	"ape/api"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	initFlag()
	checkRepo()

}

func main() {
	repo := api.SimpleRepo(*flRepo)
	apiHandler := api.NewServer(repo)
	http.Handle("/", apiHandler)
	serve()
}

func serve() {
	if *flTLS {
		serveTLS()
	}
	serveHTTP()
}

func serveHTTP() {
	port := fmt.Sprintf(":%v", *flPort)
	log.Fatal(http.ListenAndServe(port, nil))
}

func serveTLS() {
	port := fmt.Sprintf(":%v", *flPort)
	log.Fatal(http.ListenAndServeTLS(port, *flTLSCert, *flTLSKey, nil))
}

func checkRepo() {
	pkgsinfoPath := fmt.Sprintf("%v/pkgsinfo/", *flRepo)
	createDir(pkgsinfoPath)

	manifestPath := fmt.Sprintf("%v/manifests/", *flRepo)
	createDir(manifestPath)

	pkgsPath := fmt.Sprintf("%v/pkgs/", *flRepo)
	createDir(pkgsPath)

	catalogsPath := fmt.Sprintf("%v/catalogs/", *flRepo)
	createDir(catalogsPath)
}

func createDir(path string) {
	if !dirExists(path) {
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("%v must exits", path)
		}
	}
}

func dirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
