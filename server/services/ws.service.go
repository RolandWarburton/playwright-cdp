package services

import "fmt"

type PingMessage struct {
	Message   string
	SessionID string
}

var PingChannel = make(chan PingMessage)

func ListenToWSEvents(listener chan string) {
	fmt.Println("listening to events")
	// Start a goroutine to listen for events
	go func() {
		for data := range listener {
			fmt.Println("Received message:", data)
		}
	}()
}
