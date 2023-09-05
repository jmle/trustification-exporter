package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsProvider struct {
	client   *sqs.Client
	hostname string
	queue    string
}

type SqsBucket struct {
	Name string `json:"name"`
}

type SqsObject struct {
	Key string `json:"key"`
}

type SqsS3 struct {
	Object SqsObject `json:"object"`
	Bucket SqsBucket `json:"bucket"`
}

type SqsRecord struct {
	EventName string `json:"eventName"`
	S3        SqsS3  `json:"s3"`
}

type SqsMessage struct {
	Records []SqsRecord `json:"Records"`
}

func (m SqsMessage) GetEvent() EventName {
	if m.Records[0].EventName == "ObjectCreated:Put" {
		return PUT
	}
	return ""
}

func (m SqsMessage) GetBucket() string {
	return m.Records[0].S3.Bucket.Name
}

func (m SqsMessage) GetItem() string {
	return m.Records[0].S3.Object.Key
}

func NewSqsProvider(options interface{}) (SqsProvider, error) {
	sqsHostname := CheckAndReturn("SQS_HOSTNAME")
	sqsPort := CheckAndReturn("SQS_PORT")
	sqsQueue := CheckAndReturn("SQS_QUEUE")
	sqsProvider := SqsProvider{}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Println("Error loading AWS SDK config:", err)
		return SqsProvider{}, err
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s:%s", sqsHostname, sqsPort))
		o.Region = "us-east-1"
	})
	sqsProvider.client = client
	sqsProvider.hostname = sqsHostname
	sqsProvider.queue = sqsQueue

	return sqsProvider, nil
}

func (s SqsProvider) ReceiveMessage() (Message, error) {
	addr := fmt.Sprintf("http://%s:4566/000000000000/%s", s.hostname, s.queue)

	// Receive messages from the queue
	receiveInput := &sqs.ReceiveMessageInput{
		QueueUrl:            &addr,
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     10,
	}

	for {
		receiveOutput, err := s.client.ReceiveMessage(context.TODO(), receiveInput)
		if err != nil {
			fmt.Println("Error receiving message:", err)
		}

		messages := receiveOutput.Messages
		if len(messages) > 0 {
			message := receiveOutput.Messages[0]
			fmt.Println("Received message:", *message.Body)

			// Create a variable of the struct type to store the decoded data
			var msg SqsMessage

			// Unmarshal the JSON data into the struct
			err := json.Unmarshal([]byte(*message.Body), &msg)
			if err != nil {
				fmt.Println("Error:", err)
				return SqsMessage{}, err
			}

			// Delete the received message from the queue
			deleteInput := &sqs.DeleteMessageInput{
				QueueUrl:      &addr,
				ReceiptHandle: message.ReceiptHandle,
			}
			_, err = s.client.DeleteMessage(context.TODO(), deleteInput)
			if err != nil {
				fmt.Println("Error deleting message:", err)
			}
			fmt.Println("Message deleted from the queue")

			return msg, nil
		}
	}
}

func (s SqsProvider) Close() error {
	//TODO implement me
	panic("implement me")
}
