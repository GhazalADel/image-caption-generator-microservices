package main

import (
	"ImageGenerationService/handler"
	"log"
	"time"
)

func main() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := handler.ProcessReadyRequests()
			if err != nil {
				log.Printf("Error processing ready requests: %v", err)
			}
		}
	}

}
