package server

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/cfi2017/bl3-save/internal/item"
	"github.com/cfi2017/bl3-save/internal/shared"
	"github.com/cfi2017/bl3-save/pkg/pb"
	"github.com/cfi2017/bl3-save/pkg/profile"
	"github.com/gin-gonic/gin"
)

func getProfile(c *gin.Context) {
	f, err := os.Open(pwd + "/profile.sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	s, p := profile.Deserialize(f)
	c.JSON(200, struct {
		Save    shared.SavFile `json:"save"`
		Profile pb.Profile     `json:"profile"`
	}{Save: s, Profile: p})
}

func updateProfile(c *gin.Context) {
	var d struct {
		Save    shared.SavFile `json:"save"`
		Profile pb.Profile     `json:"profile"`
	}
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	backup(pwd, "profile")
	f, err := os.Create(pwd + "/profile.sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	profile.Serialize(f, d.Save, d.Profile)
	c.Status(204)
	return

}

func getBankRequest(c *gin.Context) {
	f, err := os.Open(pwd + "/profile.sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	_, p := profile.Deserialize(f)
	items := make([]item.Item, 0)
	for _, data := range p.BankInventoryList {
		d := make([]byte, len(data))
		copy(d, data)
		i, err := item.Deserialize(d)
		if err != nil {
			log.Println(err)
			log.Println(base64.StdEncoding.EncodeToString(data))
			// c.AbortWithStatus(500)
			// return
		}
		i.Wrapper = &pb.OakInventoryItemSaveGameData{
			ItemSerialNumber: data,
		}
		items = append(items, i)
	}
	c.JSON(200, &items)
	return
}

func updateBankRequest(c *gin.Context) {
	f, err := os.Open(pwd + "/profile.sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	err = f.Close()
	s, p := profile.Deserialize(f)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	var items []item.Item
	err = c.BindJSON(&items)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	backup(pwd, "profile")
	pba, err := itemsToPBArray(items)
	p.BankInventoryList = make([][]byte, len(pba))
	for i := range pba {
		p.BankInventoryList[i] = pba[i].ItemSerialNumber
	}
	f, err = os.Create(pwd + "/profile.sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	profile.Serialize(f, s, p)
	c.Status(204)
	return

}
