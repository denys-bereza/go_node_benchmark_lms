package main

import (
	"fmt"
	"net/http"
	"runtime"

	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func getPresigned(s3Client *s3.S3, path string) string {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(""),
		Key:    aws.String(path),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Println("Failed to sign request", err)
	}

	return urlStr
}

func main() {
	fmt.Printf("GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Static file server
	router.Static("/courses", "./courses")

	// JSON
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// Presigned
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-central-1"),
		Credentials: credentials.NewStaticCredentials("", "", "")},
	)
	s3Client := s3.New(sess)
	router.GET("/presigned/*resource", func(c *gin.Context) {
		name := c.Param("resource")
		link := getPresigned(s3Client, name)
		c.Redirect(http.StatusTemporaryRedirect, link)
	})
	fmt.Printf("listen and serve on 0.0.0.0:8080")

	router.Run() // listen and serve on 0.0.0.0:8080
}
