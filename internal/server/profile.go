package server

import (
	"os"

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
