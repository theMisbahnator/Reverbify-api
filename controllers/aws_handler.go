package controllers

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func upload(path string, fileName string) {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(path)
	if err != nil {
		log.Println("Error opening file:", err)
		return
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
		return
	}

	log.Println("Successfully uploaded file to", result.Location)
}
