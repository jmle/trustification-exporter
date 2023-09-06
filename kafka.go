package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"strings"
)

type KafkaProvider struct {
	// Kafka-specific configuration fields
	reader *kafka.Reader
}

type KafkaMessage struct {
	EventName string `json:"EventName"`
	Key       string `json:"Key"`
}

func (m KafkaMessage) GetEvent() EventName {
	if m.EventName == "s3:ObjectCreated:Put" {
		return PUT
	}
	return ""
}

func (m KafkaMessage) GetBucket() string {
	info := strings.Split(m.Key, "/")
	if len(info) < 2 {
		// TODO: handle error
		fmt.Printf("Skipping message with key %s: invalid format", m.Key)
	}
	return info[0]
}

func (m KafkaMessage) GetItem() string {
	info := strings.Split(m.Key, "/")
	if len(info) < 2 {
		// TODO: handle error
		fmt.Printf("Skipping message with key %s: invalid format", m.Key)
	}
	return info[1]
}

func NewKafkaProvider(mpConfig MessageProviderConfig) (KafkaProvider, error) {
	kafkaHostname := Config.kafkaHostname
	kafkaPort := Config.kafkaPort
	kafkaTopic := mpConfig.queue

	kafkaProvider := KafkaProvider{}
	kafkaProvider.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%s", kafkaHostname, kafkaPort)},
		Topic:     kafkaTopic,
		Partition: 0,
		MaxBytes:  10e6,
	})
	kafkaProvider.reader.SetOffset(kafka.LastOffset)

	return kafkaProvider, nil
}

func (k KafkaProvider) ReceiveMessage() (Message, error) {
	m, err := k.reader.ReadMessage(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

	msg := KafkaMessage{}
	err = json.Unmarshal(m.Value, &msg)
	if err != nil {
		fmt.Println("error parsing JSON:", err)
	}

	return msg, err
}

func (k KafkaProvider) Close() error {
	if err := k.reader.Close(); err != nil {
		log.Fatal("Failed to close reader:", err)
		return err
	}

	return nil
}
