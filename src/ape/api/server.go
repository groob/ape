package api

import (
	"ape/datastore"
	"log"
	"net/http"
)

type config struct {
	db  datastore.Datastore
	mux http.Handler
}

// NewServer returns an http Handler
func NewServer(options ...func(*config) error) http.Handler {
	conf := &config{}
	for _, option := range options {
		if err := option(conf); err != nil {
			log.Fatal(err)
		}
	}
	if conf.db == nil {
		log.Fatal("No datastore configured")
	}
	conf.mux = router(conf)
	return conf.mux
}

// SimpleRepo adds a file based backend
func SimpleRepo(path string) func(*config) error {
	return func(c *config) error {
		c.db = &datastore.SimpleRepo{Path: path}
		return nil
	}
}
