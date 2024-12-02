package route

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/controller"
	"github.com/AsrofunNiam/lets-code-chatbot-query/repository"
	"github.com/AsrofunNiam/lets-code-chatbot-query/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ChatBotRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	ChatBotServices := service.NewChatBotService(
		repository.NewSchemaRepository(),
		db,
		validate,
	)
	chatBotController := controller.NewChatBotController(ChatBotServices)
	router.POST("/chatbot", chatBotController.Create)
}
