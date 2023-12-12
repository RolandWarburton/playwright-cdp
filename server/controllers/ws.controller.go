package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	services "github.com/rolandwarburton/playwright-server/services"
)

type WSController struct {
	Controller
}

type PingMessage struct {
	Message   string
	SessionID string
}

var connections = make(map[string]*websocket.Conn)

func NewWSController() *WSController {
	controller := &WSController{
		Controller: Controller{
			name: "test controller",
		},
	}
	return controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (controller *WSController) WS(c *gin.Context) {
	services.HandlePingEvents()

	fmt.Println(c.Writer.Written())
	if c.Writer.Written() {
		fmt.Println("catching already written headers 1")
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer conn.Close()

	// Store the session
	var sessionID = c.Param("id")
	connections[sessionID] = conn

	// Handle WebSocket messages
	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		// message from the client
		fmt.Println(string(msg))
		if err != nil {
			// Check if the error is "websocket: close 1001"
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				fmt.Println("Client disconnected")
				break
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("Client disconnected")
				break
			}

			// Handle other errors gracefully
			c.Error(err) // Log the error

			// Send an appropriate response to the client
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
			return
		}

		// Example JSON response
		response := gin.H{
			"message": fmt.Sprintf("Hello, World! Your ID is %s", sessionID),
		}

		// Write response to WebSocket
		// if the network is bad things can go wrong here and an error will be thrown
		// however we cannot write to the socket or return JSON because
		// the connection is broken so we just return
		err = conn.WriteJSON(response)
		if err != nil {
			return
		}

		if string(msg) == "ping" {
			pingMessage := services.PingMessage{
				Message:   string(msg),
				SessionID: sessionID,
			}
			services.PingChannel <- pingMessage
		}
	}
}
