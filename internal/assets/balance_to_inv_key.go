package assets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	btik     map[string]string
	btikOnce = sync.Once{}
)

func GetBtik() map[string]string {
	var err error
	btikOnce.Do(func() {
		_ = os.MkdirAll("assets", os.ModePerm)
		// err = downloadAsset("assets/balance_to_inv_key.json", fmt.Sprintf("%s/assets/balance_to_inv_key.json", publisher))
		if err != nil {
			log.Println("couldn't download fresh assets")
		}
		btik, err = loadPartMap("assets/balance_to_inv_key.json")
	})
	if err != nil {
		panic(err)
	}
	return btik
}

func loadPartMap(file string) (m map[string]string, err error) {
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, &m)
	return
}
