package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"time"
)

func DownloadFile(hostname string, port string, bucket string, item string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		return
	}

	addr := fmt.Sprintf("http://%s:%s/%s/", hostname, port, bucket)
	cfg.Region = "us-east-1" // TODO: custom region?

	// Create an S3 client
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(addr)
	})

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the local file to save the downloaded content
	file, err := os.Create(item)
	if err != nil {
		fmt.Printf("Error creating file %s: %s", file.Name(), err)
		return
	}
	defer file.Close()

	downloader := manager.NewDownloader(client)
	numBytes, err := downloader.Download(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	if err != nil {
		fmt.Printf("Error downloading file %s: %s\n", file.Name(), err)
		return
	}

	fmt.Printf("File downloaded successfully! Downloaded %d bytes\n", numBytes)
}
