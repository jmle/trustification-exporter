package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"os/exec"
	"time"
)

// TODO: use goroutines?
func main() {
	s3hostname := CheckAndReturn("S3_HOSTNAME")
	s3port := CheckAndReturn("S3_PORT")

	mp, _ := GetMessageProvider()
	for {
		m, err := mp.ReceiveMessage()
		if err != nil {
			fmt.Printf("Error while receiving message: %s\n", err)
			continue
		}
		if m.GetEvent() != PUT {
			continue
		}

		downloadFile(s3hostname, s3port, m.GetBucket(), m.GetItem())
		//ingestFile(item)
		//removeFile(item)
	}

	mp.Close()
}

func removeFile(item string) {
	err := os.Remove(item)
	if err != nil {
		fmt.Printf("Error removing downloaded file: %s", err)
		return
	}
}

func downloadFile(hostname string, port string, bucket string, item string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		return
	}

	//addr := fmt.Sprintf("http://" + hostname + ":9000/" + bucket)
	addr := fmt.Sprintf("http://%s:%s/%s/", hostname, port, bucket)
	cfg.Region = "us-east-1"

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

func ingestFile(item string) {
	cmd := exec.Command("./guacone", "collect", "files", item)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	fmt.Printf("Ingested %s", out.String())

	err = os.Remove(item)
	if err != nil {
		fmt.Printf("Error removing downloaded file: %s", err)
		return
	}
}
