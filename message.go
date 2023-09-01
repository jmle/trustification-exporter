package main

import "os"

type Message interface {
	GetBucket() string
	GetItem() string
}

type MessageProvider interface {
	ReceiveMessage() (Message, error)
	Close() error
}

func GetMessageProvider() (MessageProvider, error) {
	providerOption := os.Getenv("MESSAGE_PROVIDER")
	switch providerOption {
	case "sqs":
		return NewSqsProvider(nil)
	default:
		return NewKafkaProvider(nil)
	}
}
