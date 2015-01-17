package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const BaseURL = "http://actualfreedom.com.au/"

func curl(url string) ([]byte, error) {
	if resp, err := http.Get(url); err != nil {
		return nil, err
	} else {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err != nil {
			return nil, err
		} else {
			return body, nil
		}
	}
}

func FixFontSize(body []byte, size float64) []byte {
	// The 'size' attribute of the 'font' tag is apparently not
	// supported in HTML5. Which could explain why changing the size
	// attribute has no effect, at least in Safari.
	body = bytes.Replace(
		body,
		[]byte("font size=\"6\""),
		[]byte(fmt.Sprintf("font style=\"font-size: %vem;\"", size)),
		-1)
	return body
}

func parsePath(c *gin.Context) (float64, string) {
	zoomS := c.Params.ByName("zoom")
	path := strings.Replace(c.Request.URL.Path, "/"+zoomS, "", 1)
	zoom, err := strconv.ParseFloat(zoomS, 32)
	if err != nil {
		zoom = 1.0
	}
	return zoom, path
}

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/2.5/richard/automaticdisplay/pageone.htm")
	})
	router.GET("/:zoom/richard/images/*rest", func(c *gin.Context) {
		_, path := parsePath(c)
		c.Redirect(302, BaseURL+path)
	})
	router.GET("/:zoom/richard/automaticdisplay/*rest", func(c *gin.Context) {
		zoom, path := parsePath(c)
		url := BaseURL + path

		if body, err := curl(url); err != nil {
			c.String(500, fmt.Sprintf("ERROR: %v", err))
		} else {
			body = FixFontSize(body, zoom)
			c.Data(200, "text/html", body)
		}
	})
	router.Run("0.0.0.0:8080")
}
