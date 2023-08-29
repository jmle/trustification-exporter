package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/segmentio/kafka-go"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Message struct {
	Key string `json:"Key"`
}

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

		handleMessage(minioHostname, msg.Key)
	}

	if err := r.Close(); err != nil {
		log.Fatal("Failed to close reader:", err)
	}
}

func handleMessage(hostname string, key string) {
	fmt.Printf("Handling message with %s", key)

	info := strings.Split(key, "/")
	if len(info) < 2 {
		fmt.Printf("Skipping message with key %s: invalid format", key)
		return
	}

	bucket := info[0]
	item := info[1]

	addr := fmt.Sprintf(hostname+":9000/%s", bucket)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"), // TODO: Default minio region - make configurable
		Endpoint:    &addr,
		DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewEnvCredentials(),
	})
	if err != nil {
		exitErrorf("Error creating session: %s", err)
	}

	downloader := s3manager.NewDownloader(sess)

	file, err := os.Create(item)
	if err != nil {
		exitErrorf("Unable to open file %q, %v", item, err)
	}

	defer file.Close()

	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		exitErrorf("Unable to download item %q, %v", item, err)
	}

	fmt.Printf("Downloaded %s, %d bytes\n", file.Name(), numBytes)

	cmd := exec.Command("./guacone", "collect", "files", item)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
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

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
