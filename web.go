package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
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
			c.Data(200, "text/html", body)
		}
	})
	router.Run("0.0.0.0:8080")
}
