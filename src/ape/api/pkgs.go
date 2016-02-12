package api

import (
	"ape/datastore"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handlePkgsCreate(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		path := ps.ByName("path")
		fmt.Println(path)
		err := db.NewPkg(path, r.Body)
		if err != nil {
			log.Fatal(err)
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("done"))
		return
	}
}
