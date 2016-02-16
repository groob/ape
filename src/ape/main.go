package main

import (
	"ape/api"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	flRepo = flag.String("repo", envString("MUNKI_REPO_PATH", ""), "path to munki repo")
	flPort = flag.String("port", envString("APE_HTTP_LISTEN_PORT", ""), "port to listen on")
)

const usage = "usage: MUNKI_REPO_PATH= APE_HTTP_LISTEN_PORT= ape -repo MUNKI_REPO_PATH -port APE_HTTP_LISTEN_PORT"

func envString(key, def string) string {
	if env := os.Getenv(key); env != "" {
		return env
	}
	return def
}

func init() {
	flag.Parse()
	if *flRepo == "" {
		flag.Usage()
		log.Fatal(usage)
	}
	if *flPort == "" {
		log.Println("no port flag specified. Using port 80 by default")
		*flPort = "80"
	}
	checkRepo()

}

func main() {
	repo := api.SimpleRepo(*flRepo)
	apiHandler := api.NewServer(repo)
	http.Handle("/", apiHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *flPort), nil))
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
