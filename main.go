package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/idelic-inc/learn-golang/aws"
	"github.com/idelic-inc/learn-golang/handler"
)

func main() {
	r := gin.Default()

	uploader, err := aws.NewS3Uploader()

	if err != nil {
		log.Printf("Could not connect to S3: %s\n", err)
		os.Exit(1)
	}

	downloader, err := aws.NewS3Downloader()

	if err != nil {
		log.Printf("Could not connect to S3: %s\n", err)
		os.Exit(1)
	}

	client, err := aws.NewS3Client()

	if err != nil {
		log.Printf("Could not connect to S3: %s\n", err)
		os.Exit(1)
	}

	// Routes
	r.POST("/images/upload/:user", func(c *gin.Context) {
		handler.ImageUpload(c, uploader)
	})
	r.GET("/images/download/:user/:name", func(c *gin.Context) {
		handler.ImageDownload(c, downloader)
	})
	r.GET("/images/list/:user", func(c *gin.Context) {
		handler.ListFiles(c, client)
	})

	// Listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	r.Run()
}
