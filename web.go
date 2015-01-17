package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
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

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/richard/automaticdisplay/pageone.htm")
	})
	router.GET("/richard/images/*rest", func(c *gin.Context) {
		c.Redirect(301, BaseURL+c.Request.URL.Path)
	})
	router.GET("/richard/automaticdisplay/*rest", func(c *gin.Context) {
		url := BaseURL + c.Request.URL.Path
		if body, err := curl(url); err != nil {
			c.String(500, fmt.Sprintf("ERROR: %v", err))
		} else {
			zoomS := c.Request.URL.Query().Get("zoom")
			if zoomS == "" {
				c.Data(200, "text/html", body)
			} else {
				if zoom, err := strconv.ParseFloat(zoomS, 32); err != nil {
					c.String(404, fmt.Sprintf("Invalid zoom: %v", err))
				} else {
					body = ZoomHTMLBody(body, zoom)
					c.Data(200, "text/html", ZoomHTMLBody(body, zoom))
				}
			}
		}
	})
	router.Run("0.0.0.0:8080")
}
