package server

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"

	"github.com/cfi2017/bl3-save-core/pkg/character"
	"github.com/cfi2017/bl3-save-core/pkg/item"
	"github.com/cfi2017/bl3-save-core/pkg/pb"
	"github.com/cfi2017/bl3-save-core/pkg/shared"
	"github.com/gin-gonic/gin"
)

var (
	charPattern        = regexp.MustCompile("^([0-9a-fA-F]+)\\.sav$")
	platformsCharacter = make(Platforms)
)

func init() {
	platformsCharacter["pc"] = character.PCMagic
	platformsCharacter["ps4"] = character.PS4Magic
}

type ItemRequest struct {
	Items    []item.Item                         `json:"items"`
	Equipped []*pb.EquippedInventorySaveGameData `json:"equipped"`
	Active   []int32                             `json:"active"`
	Platform string                              `json:"platform"`
}

func listCharacters(c *gin.Context) {
	files, err := ioutil.ReadDir(pwd)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	characters := make([]CharInfo, 0)
	for _, i := range files {
		if !i.IsDir() && charPattern.MatchString(i.Name()) {
			char, err := listChar(charPattern.FindStringSubmatch(i.Name())[1])
			if err != nil {
				log.Println(err)
				c.AbortWithStatus(500)
				return
			}
			characters = append(characters, char)
		}
	}
	c.JSON(200, &characters)

}

func getSaveById(id string) (*os.File, error) {
	return os.Open(pwd + "/" + id + ".sav")
}

func getCharacterRequest(c *gin.Context) {
	id := c.Param("id")

	f, err := getSaveById(id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	s, char, platform, err := TryDeserialize(character.Deserialize, platformsCharacter, f)
	if err != nil {
		log.Printf("error deserializing save: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	// workaround for invalid json parsing values
	for _, d := range char.GbxZoneMapFodSaveGameData.LevelData {
		if d.DiscoveryPercentage > math.MaxFloat32 {
			d.DiscoveryPercentage = -1
		}
	}
	c.JSON(200, &struct {
		Save      shared.SavFile `json:"save"`
		Character pb.Character   `json:"character"`
		Platform  string         `json:"platform"`
	}{Save: s, Character: char, Platform: platform})

}

func updateCharacterRequest(c *gin.Context) {
	id := c.Param("id")

	var d struct {
		Save      shared.SavFile `json:"save"`
		Character pb.Character   `json:"character"`
	}
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	// workaround for invalid json parsing values
	for _, d := range d.Character.GbxZoneMapFodSaveGameData.LevelData {
		if d.DiscoveryPercentage == -1 {
			d.DiscoveryPercentage = math.Float32frombits(0x7F800000) // inf
		}
	}

	backup(pwd, id)
	f, err := os.Create(pwd + "/" + id + ".sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	character.Serialize(f, d.Save, d.Character, character.PCMagic)
	c.Status(204)
	return
}

type CharInfo struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Experience int32  `json:"experience"`
}

func listChar(id string) (char CharInfo, err error) {
	f, err := getSaveById(id)
	char.ID = id
	if err != nil {
		return
	}
	defer f.Close()
	_, c, _, err := TryDeserialize(character.Deserialize, platformsCharacter, f)
	if err != nil {
		return
	}
	char.Name = c.PreferredCharacterName
	char.Experience = c.ExperiencePoints
	return
}

func getItemsRequest(c *gin.Context) {
	id := c.Param("id")
	f, err := getSaveById(id)
	if err != nil {
		c.AbortWithStatus(500)
	}
	_, char, platform, err := TryDeserialize(character.Deserialize, platformsCharacter, f)
	if err != nil {
		log.Printf("error deserializing save: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	items := make([]item.Item, 0)
	for _, data := range char.InventoryItems {
		d := make([]byte, len(data.ItemSerialNumber))
		copy(d, data.ItemSerialNumber)
		i, err := item.Deserialize(d)
		if err != nil {
			log.Println(err)
			log.Println(base64.StdEncoding.EncodeToString(data.ItemSerialNumber))
			// c.AbortWithStatus(500)
			// return
		}
		i.Wrapper = data
		items = append(items, i)
	}
	ir := ItemRequest{
		Items:    items,
		Equipped: char.EquippedInventoryList,
		Active:   char.ActiveWeaponList,
		Platform: platform,
	}
	c.JSON(200, &ir)
	return
}

func updateItemsRequest(c *gin.Context) {

	id := c.Param("id")
	f, err := getSaveById(id)
	if err != nil {
		log.Printf("error getting save: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	s, char, err := character.Deserialize(f, character.PCMagic)
	if err != nil {
		log.Printf("error deserializing save: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	err = f.Close()
	if err != nil {
		log.Printf("error deserializing save: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	var ir ItemRequest
	err = c.BindJSON(&ir)
	if err != nil {
		log.Printf("error deserializing request json: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	backup(pwd, id)
	char.InventoryItems, err = itemsToPBArray(ir.Items)
	if err != nil {
		log.Printf("error converting items to save format: %v", err)
		c.AbortWithStatusJSON(500, &err)
		return
	}
	char.ActiveWeaponList = ir.Active
	char.EquippedInventoryList = ir.Equipped
	f, err = os.Create(pwd + "/" + id + ".sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	character.Serialize(f, s, char, character.PCMagic)
	c.Status(204)
	return

}

func itemsToPBArray(items []item.Item) ([]*pb.OakInventoryItemSaveGameData, error) {
	result := make([]*pb.OakInventoryItemSaveGameData, len(items))
	for index, i := range items {
		result[index] = i.Wrapper
		seed, err := item.GetSeedFromSerial(i.Wrapper.ItemSerialNumber)
		if err != nil {
			// set seed to be 0
			seed = 0
		}
		if i.Balance == "" {
			// sanity check, if the balance is empty, just write the original item back
			continue
		}
		result[index].ItemSerialNumber, err = item.Serialize(i, seed)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}
