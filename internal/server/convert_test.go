package server

import (
	"encoding/base64"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/cfi2017/bl3-save-core/pkg/assets"
	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	assets2 "github.com/cfi2017/bl3-save/internal/assets"
)

const (
	code       = "BL3(AwAAAACQ/IA5VpSBEOCjjgcksmA0JBQNp7RKQaFQKBQKhUKh0Eaj0Wg0Go1Go9FoNBqNRqPRaDQajUaj0Wg0Go1GQ+5hhAAAAAAAAAAAAA==)"
	iterations = 1000
)

func TestConvert(t *testing.T) {

	// setup
	assets.DefaultAssetLoader = assets2.HttpAssetsLoader{}
	text := strings.Repeat(code, iterations)
	text += "asdfasdf asdf\n some random stuff"
	codes, err := extractBL3Codes(text)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()

	items := make([]item.Item, len(codes))
	for index, code := range codes {
		bs, err := base64.StdEncoding.DecodeString(code)
		if err != nil {
			log.Fatal(err)
		}
		i, err := item.Deserialize(bs)
		if err != nil {
			t.Fatal(err)
		}
		i.Wrapper = &pb.OakInventoryItemSaveGameData{
			ItemSerialNumber:    bs,
			PickupOrderIndex:    200,
			Flags:               3,
			WeaponSkinPath:      "",
			DevelopmentSaveData: nil,
		}
		items[index] = i
	}

	t.Logf("Deserialized %d items in %v", len(items), time.Since(start))

}
