package main

import (
	"fmt"
	"os"
)

type Configuration struct {
	messageProvider string
	queues          string
	s3Hostname      string
	s3Port          string
	sqsHostname     string
	sqsPort         string
	kafkaHostname   string
	kafkaPort       string
}

var Config Configuration

func InitConfiguration() {
	Config = Configuration{
		messageProvider: getEnvVar("MESSAGE_PROVIDER", false),
		queues:          getEnvVar("QUEUES", true),
		s3Hostname:      getEnvVar("S3_HOSTNAME", true),
		s3Port:          getEnvVar("S3_PORT", true),
		sqsHostname:     getEnvVar("SQS_HOSTNAME", false),
		sqsPort:         getEnvVar("SQS_PORT", false),
		kafkaHostname:   getEnvVar("KAFKA_HOSTNAME", false),
		kafkaPort:       getEnvVar("KAFKA_PORT", false),
	}
}

func getEnvVar(varName string, mandatory bool) string {
	environmentVariable := os.Getenv(varName)
	if mandatory && environmentVariable == "" {
		fmt.Println(fmt.Errorf("environment variable %s not found", varName))
		os.Exit(1)
	}
	return environmentVariable
}
