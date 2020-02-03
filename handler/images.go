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

	user := c.Param("user")
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
			err := uploadFile(uploader, file, user)
	
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

func uploadFile(uploader *s3manager.Uploader, fileHeader *multipart.FileHeader, user string) error {

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	bucketName := "golang-learn"
	keyName := user + "/" + fileHeader.Filename

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

	user := c.Param("user")
	name := c.Param("name")
	file, err := downloadFile(downloader, name, user)
	if err != nil {
		c.JSON(500, gin.H{
			"response": fmt.Sprintf("Cannot download file: %s.", err),
		})
	} else {
		c.File(file.Name())
	}
}

func downloadFile(downloader *s3manager.Downloader, fileName string, user string) (*os.File, error) {
	
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	} else {

		bucketName := "golang-learn"
	
		downParams := &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    aws.String(user + "/" + fileName),
		}
	
		_, err = downloader.Download(file, downParams)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
}

func ListFiles(c *gin.Context, client *s3.S3) {
	
	user := c.Param("user")
	filter := c.Query("filter")

	bucketName := "golang-learn"
	listParams := &s3.ListObjectsV2Input{
		Bucket: &bucketName,
		Prefix:    aws.String(user),
	}
	list, err := client.ListObjectsV2(listParams)

	var fileNames []*string
	for i := 0; i < len(list.Contents); i++ {
		key := list.Contents[i].Key
		fileName := strings.Replace(*key, user + "/", "", 1)
		if strings.Contains(fileName, filter) {
			fileNames = append(fileNames, &fileName)
		}
	}

	if err != nil {
		c.JSON(400, gin.H{
			"response": fmt.Sprintf("Cannot get files: %s.", err),
		})
	} else {
		c.JSON(200, gin.H{
			"response": fileNames,
		})
	}
}
