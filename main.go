package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/groob/ape/api"
)

func init() {
	initFlag()
	checkRepo()

}

func main() {
	var opts api.ServerOptions
	// add jwt auth
	if *flJWT {
		jwtOpt := api.JWTAuth(*flJWTSecret)
		opts = append(opts, jwtOpt)
	}
	// add basic auth
	if *flBasic {
		opts = append(opts, api.BasicAuth())
	}
	// configure repo
	repo := api.SimpleRepo(*flRepo)
	opts = append(opts, repo)
	// create handler
	apiHandler := api.NewServer(opts...)
	http.Handle("/", apiHandler)
	// serve http or https
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
