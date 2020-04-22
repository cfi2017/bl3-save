package server

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/gin-gonic/gin"
)

func convertItem(c *gin.Context) {

	var request struct {
		Base64 string `json:"base64"`
	}
	request.Base64 = strings.TrimPrefix(request.Base64, "bl3(")
	request.Base64 = strings.TrimSuffix(request.Base64, ")")
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	bs, err := base64.StdEncoding.DecodeString(request.Base64)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	var dmi item.DigitalMarineItem
	err = json.Unmarshal(bs, &dmi)
	if err != nil {
		// try deserializing item
		i, err := item.Deserialize(bs)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
			return
		}
		i.Wrapper = &pb.OakInventoryItemSaveGameData{
			ItemSerialNumber:    bs,
			PickupOrderIndex:    200,
			Flags:               3,
			WeaponSkinPath:      "",
			DevelopmentSaveData: nil,
		}
		c.JSON(200, &i)
		return
	}
	i := item.DmToGibbed(dmi)
	bs, err = item.Serialize(i, 0) // encrypt with 0 seed
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	i.Wrapper = &pb.OakInventoryItemSaveGameData{
		ItemSerialNumber:    bs,
		PickupOrderIndex:    200,
		Flags:               3,
		WeaponSkinPath:      "",
		DevelopmentSaveData: nil,
	}
	c.JSON(200, &i)
	return
}
