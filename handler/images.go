package handler

import (
	"os"
	"fmt"
	"strings"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func ImageUpload(c *gin.Context, uploader *s3manager.Uploader) {

	file, err := c.FormFile("file")
	splitstr := strings.Split(file.Filename, ".")
	if splitstr[len(splitstr) - 1] != "png" {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Cannot upload file: %s, .png file expected", file.Filename),
		})
	} else if file.Size > 1000000 {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("Cannot upload file: %s, filesize must be under 1mb", file.Filename),
		})
	} else if err != nil {
			c.JSON(400, gin.H{
				"error": fmt.Sprintf("Cannot upload file: %s", err),
			})
		} else {
			err := uploadFile(uploader, file)
	
			if err != nil {
				c.JSON(400, gin.H{
					"response": fmt.Sprintf("Cannot upload file: %s.", err),
				})
			} else {
				c.JSON(200, gin.H{
					"response": fmt.Sprintf("File uploaded successfully: %s.", file.Filename),
				})
			}
		}
}

func uploadFile(uploader *s3manager.Uploader, fileHeader *multipart.FileHeader) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "golang-learn"
	keyName := fileHeader.Filename

	upParams := &s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &keyName,
		Body:   file,
	}

	_, err = uploader.Upload(upParams)
	if err != nil {
		return err
	}

	return nil
}

func ImageDownload(c *gin.Context, downloader *s3manager.Downloader) {

	name := c.Param("name")
	
	file, err := os.Create(name)
	if err != nil {
		c.JSON(400, gin.H{
			"response": fmt.Sprintf("Cannot download file: %s.", err),
		})
	} else {

		bucketName := "golang-learn"
	
		downParams := &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    aws.String(name),
		}
	
		_, err = downloader.Download(file, downParams)
		
		if err != nil {
			c.JSON(400, gin.H{
				"response": fmt.Sprintf("Cannot upload file: %s.", err),
			})
		} else {
			c.File(file.Name())
		}

	}
}