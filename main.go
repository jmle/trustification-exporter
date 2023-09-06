package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	s3hostname := CheckAndReturn("S3_HOSTNAME")
	s3port := CheckAndReturn("S3_PORT")
	queues := strings.Split(CheckAndReturn("QUEUES"), ",")

	var wg sync.WaitGroup
	for _, queue := range queues {
		wg.Add(1)

		go func(queue string) {
			mp, _ := GetMessageProvider(MessageProviderConfig{queue: queue})
			for {
				m, err := mp.ReceiveMessage()
				if err != nil {
					fmt.Printf("Error while receiving message: %s\n", err)
					continue
				}
				if m.GetEvent() != PUT {
					continue
				}

				DownloadFile(s3hostname, s3port, m.GetBucket(), m.GetItem())
				//ingestFile(item)
				//removeFile(m.GetItem())
			}

			mp.Close()
		}(queue)
	}

	wg.Wait()
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
