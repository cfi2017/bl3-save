package assets

import (
	"log"
	"os"
	"sync"

	"github.com/cfi2017/bl3-save-core/pkg/assets"
)

var (
	db     assets.PartsDatabase
	dbOnce = sync.Once{}
)

func GetDB() assets.PartsDatabase {
	var err error
	dbOnce.Do(func() {
		_ = os.MkdirAll("assets", os.ModePerm)
		// err = downloadAsset("assets/inventory_raw.json", fmt.Sprintf("%s/assets/inventory_raw.json", publisher))
		if err != nil {
			log.Println("couldn't download fresh assets")
		}
		db, err = assets.LoadPartsDatabase("assets/inventory_raw.json")
	})
	if err != nil {
		panic(err)
	}
	return db
}
