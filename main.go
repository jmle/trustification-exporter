package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// TODO: check error handling
func main() {
	InitConfiguration()

	queues := strings.Split(Config.queues, ",")

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

				DownloadFile(Config.s3Hostname, Config.s3Port, m.GetBucket(), m.GetItem())
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
