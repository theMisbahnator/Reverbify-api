package main

import (
	ap "github.com/theMisbahnator/Reverbify/controllers"

	"github.com/gin-gonic/gin"
)

// import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", ap.Init)
	r.Run() // listen and serve on 0.0.0.0:8080
}
