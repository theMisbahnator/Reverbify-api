package controllers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func Init_get_url(c *gin.Context) {
	var body signed_url_request
	err := c.BindJSON(&body)
	if handleError(err, c, "Invalid request body") {
		return
	}

	sess := session.Must(session.NewSession())

	url, err := getSignedUrl(body.Filename, sess)
	if handleError(err, c, url) {
		return
	}

	sendUrlResponse(c, url)
}

func upload(path string, fileName string) (string, error) {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file:", err)
		return "Error opening file", err
	}
	defer file.Close()

	// Upload the file to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("reverbify"),
		Key:    aws.String(fileName),
		Body:   file,
	})

	if err != nil {
		log.Println("Error uploading file:", err)
		return "Error uploading file", err
	}

	log.Println("Successfully uploaded file to", result.Location)

	// get signed url
	url, err := getSignedUrl(fileName, sess)
	if err != nil {
		log.Println("Error getting signed url:", err)
		return "Error getting signed url", err
	}

	return url, nil
}

func getSignedUrl(fileName string, sess *session.Session) (string, error) {
	// Create a new S3 service client
	s3svc := s3.New(sess)
	bucket := "reverbify"

	// set expiration time to around 1 week
	const secondsPerDay = 24 * 60 * 59
	const oneWeekInSeconds = secondsPerDay * 7
	expiration := time.Duration(oneWeekInSeconds) * time.Second

	// Generate the signed URL
	req, _ := s3svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	url, err := req.Presign(expiration)
	if err != nil {
		fmt.Println("Error signing URL:", err)
		return "Error signing URL", err
	}

	return url, nil
}
