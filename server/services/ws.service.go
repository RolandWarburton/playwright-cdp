package services

import "fmt"

// "github.com/gin-gonic/gin"
// "github.com/rolandwarburton/playwright-server/errors"
// database "github.com/rolandwarburton/playwright-server/models"
// "gorm.io/gorm"

type PingMessage struct {
	Message   string
	SessionID string
}

var PingChannel = make(chan PingMessage)

func HandlePingEvents() {
	for {
		message := <-PingChannel
		fmt.Printf("Received ping from session ID: %s '%s'\n", message.SessionID, message.Message)
		// Handle the ping event here
	}
}

// type IWSService interface {
// 	handlePingEvents() (gin.H, *errors.RestError)
// }
