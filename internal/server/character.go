package server

import (
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/cfi2017/bl3-save/internal/shared"
	"github.com/cfi2017/bl3-save/pkg/character"
	"github.com/cfi2017/bl3-save/pkg/pb"
	"github.com/gin-gonic/gin"
)

var (
	charPattern = regexp.MustCompile("(\\d+)\\.sav")
)

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
				c.AbortWithStatus(500)
				return
			}
			characters = append(characters, char)
		}
	}
	c.JSON(200, &characters)

}

func getCharacter(c *gin.Context) {
	id := c.Param("id")

	f, err := os.Open(pwd + "/" + id + ".sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	s, char := character.Deserialize(f)

	// workaround for invalid json parsing values
	for _, d := range char.GbxZoneMapFodSaveGameData.LevelData {
		if *d.DiscoveryPercentage > math.MaxFloat32 {
			*d.DiscoveryPercentage = -1
		}
	}

	c.JSON(200, &struct {
		Save      shared.SavFile `json:"save"`
		Character pb.Character   `json:"character"`
	}{Save: s, Character: char})

}

func updateCharacter(c *gin.Context) {
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
		if *d.DiscoveryPercentage == -1 {
			*d.DiscoveryPercentage = math.Float32frombits(0x7F800000) // inf
		}
	}
	f, err := os.Create(pwd + "/" + id + ".sav")
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	defer f.Close()
	character.Serialize(f, d.Save, d.Character)
	c.Status(204)
	return
}

type CharInfo struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Experience int32  `json:"experience"`
}

func listChar(id string) (char CharInfo, err error) {
	char.ID, err = strconv.Atoi(id)
	if err != nil {
		return
	}
	f, err := os.Open(pwd + "/" + id + ".sav")
	if err != nil {
		return
	}
	defer f.Close()
	_, c := character.Deserialize(f)
	char.Name = *c.PreferredCharacterName
	char.Experience = *c.ExperiencePoints
	return
}
