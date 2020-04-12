package item

import (
	"encoding/base64"
	"log"
	"testing"
)

var checks = []string{
	"A6cRHH+sfCuWGEZz2Lc5FWDbSfcQLmbaOV6SzgYP",
}

func TestDecryptSerial(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			panic(err)
		}
		item, err := DecryptSerial(bs)
		if err != nil {
			panic(err)
		}
		log.Println(item)
	}
}

func TestDeserialize(t *testing.T) {
	for _, check := range checks {
		bs, err := base64.StdEncoding.DecodeString(check)
		if err != nil {
			panic(err)
		}
		item, err := Deserialize(bs)
		if err != nil {
			panic(err)
		}
		log.Println(item)
	}
}
