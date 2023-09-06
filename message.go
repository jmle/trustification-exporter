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

type MessageProviderConfig struct {
	queue string
}

// GetMessageProvider Returns a MessageProvider with the given config. Defaults to Kafka provider if no MESSAGE_PROVIDER environment variable is found
func GetMessageProvider(config MessageProviderConfig) (MessageProvider, error) {
	providerOption := os.Getenv("MESSAGE_PROVIDER")
	switch providerOption {
	case "sqs":
		return NewSqsProvider(config)
	default:
		return NewKafkaProvider(config)
	}
}
