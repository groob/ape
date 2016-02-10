package api

import (
	"ape/datastore"
	"ape/models"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handleManifestsList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		manifests, err := db.AllManifests()
		if err != nil {
			log.Fatal(err)
		}
		view := models.ManifestList(manifests)
		jsn, err := view.ToJSON()
		if err != nil {
			log.Fatal(err)
		}
		rw.Write(jsn)
		return
	}
}

func handleManifestsShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		name := ps.ByName("name")
		manifest, err := db.Manifest(name)
		if err != nil {
			log.Println(err)
		}
		view := manifest.View()
		jsn, err := view.ToJSON()
		if err != nil {
			log.Println(err)
		}
		rw.Write(jsn)
		return
	}
}

func handleManifestsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		f := &struct {
			Name string `json:"name"`
		}{}
		err := json.NewDecoder(r.Body).Decode(f)
		switch err {
		case nil:
			break
		case io.EOF:
			rw.WriteHeader(http.StatusNoContent)
			return
		default:
			log.Println(err)
		}

		manifest, err := db.NewManifest(f.Name)
		if err != nil {
			log.Println(err)
			return
		}
		if err := db.SaveManifest(manifest); err != nil {
			log.Println(err)
			return
		}

		return
	}
}
