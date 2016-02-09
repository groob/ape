package main

import (
	"ape/api"
	"log"
	"net/http"
)

func main() {
	repo := api.GitRepo("repo/pkgsinfo")
	apiHandler := api.NewServer(repo)
	http.Handle("/", apiHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
