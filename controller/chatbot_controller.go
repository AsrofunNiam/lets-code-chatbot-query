package controller

import (
	"github.com/gin-gonic/gin"
)

type ChatBotController interface {
	Create(context *gin.Context)
}
