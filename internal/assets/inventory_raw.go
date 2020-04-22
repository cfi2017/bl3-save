package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

var (
	db     PartsDatabase
	dbOnce = sync.Once{}
)

func GetDB() PartsDatabase {
	var err error
	dbOnce.Do(func() {
		_ = os.MkdirAll("assets", os.ModePerm)
		err = downloadAsset("assets/inventory_raw.json", fmt.Sprintf("%s/assets/inventory_raw.json", publisher))
		if err != nil {
			return
		}
		db, err = loadPartsDatabase("assets/inventory_raw.json")
	})
	if err != nil {
		panic(err)
	}
	return db
}

func loadPartsDatabase(file string) (db PartsDatabase, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &db)
	return
}
