package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func GenerateGitHub(c *gin.Context, name string, width, margin int) {
	if width < 32 || width-margin < 10 {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Invalid parameters"))
		return
	}

	img := GenerateIdenticon([]byte(name), width, margin)
	c.Header("Content-Type", "image/png")
	c.Stream(func(w io.Writer) bool {
		png.Encode(w, img)
		return false
	})
}

func Generate8bit(c *gin.Context) {
	name := c.Param("name")
	gender := c.Param("gender")

	switch {
	case gender == "m":
		gender = "male"
	case gender == "f":
		gender = "female"
	case gender == "male" || gender == "female":
		//do nothing
	default:
		c.AbortWithError(http.StatusInternalServerError, errors.New("Invalid parameters"))
		return
	}

	log.Println(name)
	InitAssets()
	img := GenerateIdenticon8bits(gender, []byte(name))
	c.Header("Content-Type", "image/png")
	c.Stream(func(w io.Writer) bool {
		png.Encode(w, img)
		return false
	})
}

func main() {
	r := gin.Default()

	g := r.Group("/i", func(c *gin.Context) {
		c.Header("Cache-Control", "max-age=315360000")
	})

	g.GET("/g/:name", func(c *gin.Context) {
		name := c.Param("name")

		queryWidth := c.DefaultQuery("w", "512")
		queryMargin := c.DefaultQuery("m", "32")

		width, err := strconv.Atoi(queryWidth)
		margin, err2 := strconv.Atoi(queryMargin)

		if err != nil || err2 != nil {
			log.Println(err, err2)
			c.AbortWithError(http.StatusInternalServerError, errors.New("Invalid parameters"))
			return
		}

		GenerateGitHub(c, name, width, margin)
	})

	g.GET("/8/:gender/:name", Generate8bit)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
