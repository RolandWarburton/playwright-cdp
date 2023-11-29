package controllers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/user"

	"github.com/gin-gonic/gin"
	// "github.com/rolandwarburton/playwright-server/errors"
	// database "github.com/rolandwarburton/playwright-server/models"
)

type AgentController struct {
	Controller
}

type Message struct {
	Action string `json:"action"`
}

func NewAgentController() *AgentController {
	controller := &AgentController{
		Controller: Controller{
			name: "test controller",
		},
	}
	return controller
}

// COMMON
func (controller *AgentController) callAgentAPI(c *gin.Context, action string) {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Error:", err)
		c.AbortWithStatusJSON(500, err)
		return
	}

	uid := currentUser.Uid
	fmt.Println("User ID:", uid)
	message := Message{
		Action: action,
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	conn, err := net.Dial("unix", fmt.Sprintf("/var/run/user/%s/playwright-agent.socket", uid))
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(jsonData)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}

	result := gin.H{
		"message": "ok",
	}

	c.IndentedJSON(http.StatusOK, result)
}

func (controller *AgentController) CreateAgent(c *gin.Context) {
	controller.callAgentAPI(c, "connect")
}

func (controller *AgentController) ExampleAction(c *gin.Context) {
	controller.callAgentAPI(c, "action-example")
}
