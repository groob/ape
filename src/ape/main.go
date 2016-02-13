package main

import (
	"ape/api"
	"log"
	"net/http"
)

func main() {
	// db := &datastore.SimpleRepo{
	// 	Path: "repo",
	// }
	// pkgsinfos, err := db.AllPkgsinfos()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, pkgsinfo := range *pkgsinfos {
	// 	fmt.Println(pkgsinfo.Filename)
	// }
	// os.Exit(0)
	repo := api.SimpleRepo("repo")
	apiHandler := api.NewServer(repo)
	http.Handle("/", apiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
