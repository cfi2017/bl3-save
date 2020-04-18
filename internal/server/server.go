package server

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	pwd string
)

func init() {
	pwd, _ = os.Getwd()
}

func Start(opts Options) error {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.LoggerWithWriter(os.Stderr, "/stats"), gin.Recovery())
	cfg := cors.Config{
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	if opts.Insecure {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = []string{"https://bl3.swiss.dev", "http://localhost:4200"}
	}
	r.Use(cors.New(cfg))

	r.GET("/stats", func(c *gin.Context) {
		_, err := os.Stat(pwd + "/profile.sav")
		c.JSON(200, &struct {
			Pwd        string `json:"pwd"`
			HasProfile bool   `json:"hasProfile"`
		}{
			Pwd:        pwd,
			HasProfile: err == nil && !os.IsNotExist(err),
		})
	})

	r.POST("/cd", func(c *gin.Context) {
		var body struct {
			Pwd string `json:"pwd" binding:"required"`
		}
		err := c.Bind(&body)
		if err != nil {
			return
		}
		pwd = strings.TrimSuffix(body.Pwd, "/")
		c.JSON(200, struct {
			Pwd string `json:"pwd"`
		}{Pwd: pwd})
	})

	r.GET("/profile", getProfile)
	r.POST("/profile", updateProfile)

	r.GET("/characters", listCharacters)
	r.GET("/characters/:id", getCharacterRequest)
	r.POST("/characters/:id", updateCharacterRequest)

	r.GET("/characters/:id/items", getItemsRequest)
	r.POST("/characters/:id/items", updateItemsRequest)

	r.GET("/convert", convertItem)

	return r.Run(":5050")
}
