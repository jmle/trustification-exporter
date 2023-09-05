package main

import (
	"fmt"
	"os"
)

func CheckAndReturn(varName string) string {
	environmentVariable := os.Getenv(varName)
	if environmentVariable == "" {
		fmt.Println(fmt.Errorf("environment variable %s not found", varName))
		os.Exit(1)
	}
	return environmentVariable
}
