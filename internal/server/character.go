package server

import (
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

	"github.com/cfi2017/bl3-save/pkg/character"
	"github.com/gin-gonic/gin"
)

var (
	charPattern = regexp.MustCompile("(\\d+)\\.sav")
)

func ListCharacters(c *gin.Context) {
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
