package services

import (
	"encoding/json"
	"fmt"
	errors "github.com/rolandwarburton/playwright-server/errors"
	"net"
	"os/user"
)

type Message struct {
	Action string `json:"action"`
}

func CallAgentAPI(action string) ([]byte, *errors.HTTPError) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		return nil, &errors.HTTPError{
			Message: err.Error(),
			Status:  500,
			Error:   "Internal Server Error",
		}
	}

	uid := currentUser.Uid

	// the message from the websocket will be written to the agent socket
	message := Message{
		Action: action,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return nil, &errors.HTTPError{
			Message: err.Error(),
			Status:  500,
			Error:   "Internal Server Error",
		}
	}

	conn, err := net.Dial("unix", fmt.Sprintf("/var/run/user/%s/playwright-agent.socket", uid))
	if err != nil {
		return nil, &errors.HTTPError{
			Message: err.Error(),
			Status:  500,
			Error:   "Internal Server Error",
		}
	}
	defer conn.Close()

	_, err = conn.Write(jsonData)
	if err != nil {
		return nil, &errors.HTTPError{
			Message: err.Error(),
			Status:  500,
			Error:   "Internal Server Error",
		}
	}

	// Read the response from the Unix socket
	response := make([]byte, 1024) // Adjust the buffer size as per your requirement
	bytesRead, err := conn.Read(response)
	if err != nil {
		return nil, &errors.HTTPError{
			Message: err.Error(),
			Status:  500,
			Error:   "Internal Server Error",
		}
	}

	return response[0:bytesRead], nil
}

// func CallAgentAPI(c *gin.Context, action string) {
// 	currentUser, err := user.Current()
// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		c.AbortWithStatusJSON(500, err)
// 		return
// 	}
//
// 	uid := currentUser.Uid
// 	fmt.Println("User ID:", uid)
// 	message := Message{
// 		Action: action,
// 	}
// 	jsonData, err := json.Marshal(message)
// 	if err != nil {
// 		c.AbortWithStatusJSON(500, err)
// 		return
// 	}
//
// 	conn, err := net.Dial("unix", fmt.Sprintf("/var/run/user/%s/playwright-agent.socket", uid))
// 	if err != nil {
// 		c.AbortWithStatusJSON(500, err)
// 		return
// 	}
// 	defer conn.Close()
//
// 	_, err = conn.Write(jsonData)
// 	if err != nil {
// 		c.AbortWithStatusJSON(500, err)
// 		return
// 	}
//
// 	result := gin.H{
// 		"message": "ok",
// 	}
//
// 	c.IndentedJSON(http.StatusOK, result)
// }
