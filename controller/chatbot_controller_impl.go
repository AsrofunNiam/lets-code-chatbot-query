package controller

import (
	"net/http"

	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/AsrofunNiam/lets-code-chatbot-query/service"
	"github.com/gin-gonic/gin"
)

type ChatBotControllerImpl struct {
	ChatBotService service.ChatBotService
}

func NewChatBotController(chatBotService service.ChatBotService) ChatBotController {
	return &ChatBotControllerImpl{
		ChatBotService: chatBotService,
	}
}

func (controller *ChatBotControllerImpl) Create(c *gin.Context) {
	request := &web.ChatBotCreateRequest{}
	helper.ReadFromRequestBody(c, &request)

	chatBotResponse := controller.ChatBotService.Create(request, c)
	webResponse := web.WebResponse{
		Success: true,
		Message: chatBotResponse,
	}

	c.JSON(http.StatusOK, webResponse)
}
