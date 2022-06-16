package response

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"hms/gateway/pkg/config"
	"log"
	"net/http"
)

// Created Response with created document or status only depending on user's request header "Prefer: return={representation|minimal}"
func Created(docId string, doc interface{}, c *gin.Context) {
	location, err := location(docId, c)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", location)
	c.Header("ETag", docId)

	prefer := c.Request.Header.Get("Prefer")
	if prefer == "return=representation" {
		c.JSON(http.StatusCreated, doc)
	} else {
		c.AbortWithStatus(http.StatusCreated)
	}
}

// Make document location string
func location(docId string, c *gin.Context) (location string, err error) {
	cfgRaw, ok := c.Get("cfg")
	if !ok {
		err = errors.New("not found configuration in context")
		return
	}

	cfg, ok := cfgRaw.(*config.Config)
	if !ok {
		err = fmt.Errorf("bad configuration in context: %+v", cfgRaw)
		return
	}

	path := c.Request.RequestURI
	if "/" != string(path[len(path)-1]) {
		path += "/"
	}

	location = cfg.BaseUrl + path + docId
	return
}
