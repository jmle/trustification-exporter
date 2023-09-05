package main

import "os"

type EventName string

const (
	PUT EventName = "PUT"
)

type Message interface {
	GetEvent() EventName
	GetBucket() string
	GetItem() string
}

type MessageProvider interface {
	ReceiveMessage() (Message, error)
	Close() error
}

// GetMessageProvider Returns a MessageProvider. Defaults to Kafka provider if no MESSAGE_PROVIDER environment variable is found
func GetMessageProvider() (MessageProvider, error) {
	providerOption := os.Getenv("MESSAGE_PROVIDER")
	switch providerOption {
	case "sqs":
		return NewSqsProvider(nil)
	default:
		return NewKafkaProvider(nil)
	}
}
