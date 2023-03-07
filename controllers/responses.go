package controllers

import (
	"github.com/gin-gonic/gin"
)

type audio_request struct {
	Url    string `json:"url"`
	Pitch  string `json:"pitch"`
	Bass   bool   `json:"bass"`
	Reverb bool   `json:"reverb"`
}

type audio_response struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Duration  string `json:"duration"`
	Thumbnail string `json:"thumbnail"`
}

func handleError(err error, c *gin.Context, errMSG string) bool {
	if err != nil {
		c.JSON(400, gin.H{
			"ERROR": errMSG,
		})
		return true
	}
	return false
}

func healthCheck(c *gin.Context) {
	c.JSON(200, "This endpoint works.")
}

func sendAudioResponse(c *gin.Context, title string, duration string, author string, thumbnail string) {
	response := audio_response{
		title, author, duration, thumbnail,
	}
	c.JSON(200, response)
}
