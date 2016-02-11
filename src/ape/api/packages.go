package api

import (
	"ape/datastore"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func handlePackagesList(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// 		pkgsinfos, err := db.AllPkgsInfos()
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		view := models.PkgsInfoList(pkgsinfos)
		// 		jsn, err := view.ToJSON()
		// 		if err != nil {
		// 			log.Fatal(err)
		// 		}
		// 		rw.Write(jsn)
		// 		return
	}
}

//
func handlePackagesShow(db datastore.Datastore) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// 		name := ps.ByName("name")
		// 		pkginfo, err := db.PkgInfo(name)
		// 		if err != nil {
		// 			log.Println(err)
		// 		}
		// 		view := pkginfo.View()
		// 		jsn, err := view.ToJSON()
		// 		if err != nil {
		// 			log.Println(err)
		// 		}
		// 		rw.Write(jsn)
		// 		return
	}
}
