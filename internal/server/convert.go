package server

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/gin-gonic/gin"
)

var (
	bl3CodeRegexp = regexp.MustCompile("(bl|BL)3\\(([A-Za-z0-9+/=]+)\\)")
)

func convertItem(c *gin.Context) {

	var request struct {
		Base64 string `json:"base64"`
	}
	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, err)
		return
	}
	request.Base64 = strings.TrimSpace(request.Base64)
	bs, err := base64.StdEncoding.DecodeString(request.Base64)
	if err != nil {
		// try extracting bl3 codes
		codes, err := extractBL3Codes(request.Base64)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, err)
			return
		}
		items := make([]item.Item, len(codes))
		for index, code := range codes {
			bs, err := base64.StdEncoding.DecodeString(code)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(500, err)
				return
			}
			i, err := item.Deserialize(bs)
			if err != nil {
				log.Printf("error in bl3(%s): %v", code, err)
				c.AbortWithStatusJSON(500, err)
				return
			}
			i.Wrapper = &pb.OakInventoryItemSaveGameData{
				ItemSerialNumber: bs,
				PickupOrderIndex: 200,
				Flags:            3,
				// WeaponSkinPath:      "",
				DevelopmentSaveData: nil,
			}
			items[index] = i
		}
		c.JSON(200, &items)
		return
	}
	var dmi item.DigitalMarineItem
	err = json.Unmarshal(bs, &dmi)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, err)
		return
	}
	i := item.DmToGibbed(dmi)
	bs, err = item.Serialize(i, 0) // encrypt with 0 seed
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	i.Wrapper = &pb.OakInventoryItemSaveGameData{
		ItemSerialNumber: bs,
		PickupOrderIndex: 200,
		Flags:            3,
		// WeaponSkinPath:      "",
		DevelopmentSaveData: nil,
	}
	c.JSON(200, &[]item.Item{i})
	return
}

func extractBL3Codes(text string) (codes []string, err error) {
	matches := bl3CodeRegexp.FindAllStringSubmatch(text, -1)
	codes = make([]string, len(matches))
	for i, match := range matches {
		codes[i] = match[2]
	}
	return
}
