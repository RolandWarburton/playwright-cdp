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
		// iterate over each value on the listener as they arrive
		for data := range listener {
			switch data {
			case "ping":
				fmt.Println("pong")
			default:
				fmt.Println("idk")
			}
		}
	}()
}
