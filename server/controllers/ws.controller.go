package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	services "github.com/rolandwarburton/playwright-server/services"
	playwright_util "github.com/rolandwarburton/playwright-server/utils"
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
	emitter := playwright_util.NewEmitter()

	// Create a listener channel
	listener := make(chan string)

	// Register the listener for the "message" event
	emitter.On("message", listener)

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

	// start reacting to websocket events via the service
	services.ListenToWSEvents(listener)

	// Handle WebSocket messages
	for {
		// Read message from WebSocket
		_, msg, err := conn.ReadMessage()
		// message from the client
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

			// Check if headers have already been written
			if !c.Writer.Written() {
				// Send an appropriate response to the client
				c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred"})
			}
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
			fmt.Println("Failed to write JSON response: ", err)
			break
		}

		if string(msg) != "" {
			emitter.Emit("message", string(msg))
		}
	}
}
