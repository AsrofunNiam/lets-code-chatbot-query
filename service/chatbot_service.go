package service

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
)

type ChatBotService interface {
	Create(request *web.ChatBotCreateRequest, c *gin.Context) *genai.GenerateContentResponse
}
