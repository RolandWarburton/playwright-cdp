package services

import (
	"fmt"
)

type PingMessage struct {
	Message   string
	SessionID string
}

var PingChannel = make(chan PingMessage)

type WSEvent struct {
	Response string
	Data     interface{}
}

func ListenToWSEvents(listener chan string, sender chan WSEvent) {
	fmt.Println("listening to events")
	// Start a goroutine to listen for events
	go func() {
		// iterate over each value on the listener as they arrive
		for data := range listener {
			fmt.Println(data)
			switch data {
			case "CONNECT_TO_AGENT":
				res := WSEvent{Response: "unknown"}
				err := CallAgentAPI("connect")
				if err != nil {
					fmt.Println(err)
					res.Response = "error"
				} else {
					res.Response = "ok"
				}
				sender <- res
			default:
				fmt.Println("idk")
			}
		}
	}()
}
