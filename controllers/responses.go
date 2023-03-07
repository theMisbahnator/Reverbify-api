package controllers

import (
	"github.com/gin-gonic/gin"
)

type audio_request struct {
	Url    string `json:"url"`
	Pitch  string `json:"pitch"`
	Bass   bass   `json:"bass"`
	Reverb bool   `json:"reverb"`
}

type bass struct {
	SetBass     bool   `json:"change"`
	CentFreq    string `json:"centerFreq"`
	FilterWidth string `json:"filterWidth"`
	Gain        string `json:"gain"`
}

type audio_response struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Duration  string `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	AwsUrl    string `json:"signedUrl"`
	Filename  string `json:"filename"`
}

type signed_url_request struct {
	Filename string `json:"filename"`
}

type signed_url_response struct {
	Url string `json:"signedUrl"`
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

func sendAudioResponse(c *gin.Context, title string, duration string, author string, thumbnail string,
	signedUrl string, filename string) {
	response := audio_response{
		title, author, duration, thumbnail, signedUrl, filename,
	}
	c.JSON(200, response)
}

func sendUrlResponse(c *gin.Context, url string) {
	response := signed_url_response{
		url,
	}
	c.JSON(200, response)
}
