package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Message struct {
	Key string `json:"Key"`
}

type Bus interface {
	NewListener() Bus
	Read() Message
}

// TODO: only act on put for the moment
// TODO: use goroutines?
// TODO: generic messaging bus interface (kafka/sqs)
// TODO: generic storage interface (s3/minio)
func main() {
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaHostname := os.Getenv("KAFKA_HOSTNAME")
	minioHostname := os.Getenv("MINIO_HOSTNAME")

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaHostname + ":9092"},
		Topic:     kafkaTopic,
		Partition: 0,
		MaxBytes:  10e6,
	})
	r.SetOffset(kafka.LastOffset)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		var msg Message
		err = json.Unmarshal(m.Value, &msg)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}

		info := strings.Split(msg.Key, "/")
		if len(info) < 2 {
			fmt.Printf("Skipping message with key %s: invalid format", msg.Key)
			return
		}

		bucket := info[0]
		item := info[1]

		downloadFile(minioHostname, bucket, item)
		ingestFile(item)
		removeFile(item)
	}

	if err := r.Close(); err != nil {
		log.Fatal("Failed to close reader:", err)
	}
}

func removeFile(item string) {
	err := os.Remove(item)
	if err != nil {
		fmt.Printf("Error removing downloaded file: %s", err)
		return
	}
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

func downloadFile(hostname string, bucket string, item string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		return
	}

	addr := fmt.Sprintf("http://" + hostname + ":9000/" + bucket)
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
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	downloader := manager.NewDownloader(client)
	numBytes, err := downloader.Download(ctx, file, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})

	fmt.Printf("File downloaded successfully! Downloaded %d bytes", numBytes)
}
