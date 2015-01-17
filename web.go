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

// ZoomHTMLBody transforms the HTML body such that it renders zoomed
// in the manner of browser-level zoom.
func ZoomHTMLBody(body []byte, zoom float64) []byte {
	// Courtesy of http://stackoverflow.com/a/1156526/55246
	zoomCssTmpl := "body { zoom: %v; -moz-transform: scale(%v); -moz-transform-origin: 0 0}"
	zoomCss := fmt.Sprintf(zoomCssTmpl, zoom, zoom)
	afSiteStyleTag := "</style>"
	return bytes.Replace(
		body,
		[]byte(afSiteStyleTag),
		[]byte(zoomCss+afSiteStyleTag), 1)
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
		c.Redirect(301, "/2.5/richard/automaticdisplay/pageone.htm")
	})
	router.GET("/:zoom/richard/images/*rest", func(c *gin.Context) {
		_, path := parsePath(c)
		c.Redirect(301, BaseURL+path)
	})
	router.GET("/:zoom/richard/automaticdisplay/*rest", func(c *gin.Context) {
		zoom, path := parsePath(c)
		url := BaseURL + path

		if body, err := curl(url); err != nil {
			c.String(500, fmt.Sprintf("ERROR: %v", err))
		} else {
			body = ZoomHTMLBody(body, zoom)
			c.Data(200, "text/html", body)
		}
	})
	router.Run("0.0.0.0:8080")
}
