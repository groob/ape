package datastore

import (
	"ape/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/groob/plist"
)

var makecatalogs = make(chan bool, 1)

func (r *SimpleRepo) makeCatalogs(done chan bool) {
	t1 := time.Now()
	catalogs := map[string]*models.Catalogs{}
	pkgsinfos, err := r.AllPkgsinfos()
	if err != nil {
		log.Println(err)
	}
	allCatalogs := pkgsinfos.Catalog("all")
	catalogs["all"] = allCatalogs
	for _, info := range *allCatalogs {
		for _, catalogName := range info.Catalogs {
			catalogs[catalogName] = pkgsinfos.Catalog(catalogName)
		}
	}

	for k, v := range catalogs {
		err = r.saveCatalog(k, v)
		if err != nil {
			log.Println(err)
		}
	}
	t2 := time.Now()
	log.Printf("[%s] Number of Pkgsinfos: %v time: %v\n", "makecatalogs", len(*pkgsinfos), t2.Sub(t1))
	done <- true
}

func (r *SimpleRepo) saveCatalog(name string, catalogs *models.Catalogs) error {
	catalogsPath := fmt.Sprintf("%v/catalogs/%v", r.Path, name)
	var file *os.File
	var err error
	if _, err := os.Stat(catalogsPath); err != nil {
		file, err = os.Create(catalogsPath)
	} else {
		file, err = os.OpenFile(catalogsPath, os.O_WRONLY, 0755)
	}
	if err != nil {
		return err
	}
	defer file.Close()
	enc := plist.NewEncoder(file)
	enc.Indent("  ")
	return enc.Encode(catalogs)

}

//WatchCatalogs creates catalogs from pkgsinfos
func (r *SimpleRepo) WatchCatalogs() {
	done := make(chan bool, 1)
	for {
		go r.makeCatalogs(done)
		<-done
		<-makecatalogs
	}
}
