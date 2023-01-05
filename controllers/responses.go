package controllers

import (
	"github.com/gin-gonic/gin"
)

type audio_request struct {
	Url       string `json:"url"`
	PitchType int    `json:"pitch"`
}

type audio_response struct {
	Title     string `json:"title"`
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

func sendAudioResponse(c *gin.Context, title string, duration string, thumbnail string) {
	response := audio_response{
		title, duration, thumbnail,
	}
	c.JSON(200, response)
}
