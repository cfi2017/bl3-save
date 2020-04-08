package server

import (
	"os"
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

func Start() error {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://bl3.swiss.dev", "http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, &struct {
			Pwd string `json:"pwd"`
		}{
			Pwd: pwd,
		})
	})

	r.POST("/cd", func(c *gin.Context) {
		var body struct{ path string }
		err := c.BindJSON(&body)
		if err != nil {
			return
		}
		pwd = body.path
		c.JSON(200, struct {
			Pwd string `json:"pwd"`
		}{Pwd: pwd})
	})

	r.GET("/profile", getProfile)
	r.POST("/profile", updateProfile)

	r.GET("/characters", ListCharacters)
	r.GET("/characters/:id")
	r.POST("/characters/:id")

	return r.Run(":5050")
}
