package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rolandwarburton/playwright-server/services"
)

type AgentController struct {
	Controller
}

func NewAgentController() *AgentController {
	controller := &AgentController{
		Controller: Controller{
			name: "test controller",
		},
	}
	return controller
}

func (controller *AgentController) CreateAgent(c *gin.Context) {
	services.CallAgentAPI("connect")
}

func (controller *AgentController) ExampleAction(c *gin.Context) {
	services.CallAgentAPI("action-example")
}
