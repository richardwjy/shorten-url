package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST Shorten Request - Response
type shortenReq struct {
	Link string `json:"link"`
}

type shortenRes struct {
	Link       string `json:"link,omitempty"`
	ShortenUrl string `json:"shorten_link,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode string `json:"status_code,omitempty"`
}

// GET Redirect Request - Response
type redirectReq struct {
	ShortLink string `uri:"short_link" binding:"required"`
}

const SuccessCode string = "00"
const FailedCode string = "10"

var urlMap = make(map[string]string)

func main() {
	router := gin.Default()
	router.GET("/:short_link", getRedirect)
	router.POST("/shorten", postShorten)

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}

func shorten(longUrl string) (string, error) {
	hash := sha256.Sum256([]byte(longUrl))
	shortCode := base64.URLEncoding.EncodeToString(hash[:6])
	if val, ok := urlMap[shortCode]; ok {
		if val == longUrl {
			return shortCode, nil
		}
		return shorten(longUrl)
	}
	urlMap[shortCode] = longUrl
	return shortCode, nil
}

func getRedirect(c *gin.Context) {
	var req redirectReq
	if err := c.ShouldBindUri(&req); err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(req.ShortLink) < 6 {
		c.JSON(http.StatusBadRequest, &shortenRes{
			Message:    "Invalid short link",
			StatusCode: FailedCode,
		})
	}
	if val, ok := urlMap[req.ShortLink]; ok {
		c.Redirect(http.StatusPermanentRedirect, val)
		return
	}
	c.JSON(http.StatusBadRequest, &shortenRes{
		Message:    "Short Link not found",
		StatusCode: FailedCode,
	})
}

func postShorten(c *gin.Context) {
	var req shortenReq
	var shortLink string
	var err error
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Error while binding: %s\n", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, &shortenRes{
			Message:    err.Error(),
			StatusCode: FailedCode,
		})
		return
	}
	if shortLink, err = shorten(req.Link); err != nil {
		log.Printf("Error while shortening: %s\n", err.Error())
		c.IndentedJSON(http.StatusInternalServerError, &shortenRes{
			Message:    err.Error(),
			StatusCode: FailedCode,
		})
		return
	}
	res := &shortenRes{
		Link:       req.Link,
		ShortenUrl: shortLink,
		StatusCode: SuccessCode,
	}
	c.IndentedJSON(http.StatusCreated, res)
}
