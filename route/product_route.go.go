package route

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/controller"
	"github.com/AsrofunNiam/lets-code-chatbot-query/repository"
	"github.com/AsrofunNiam/lets-code-chatbot-query/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ProductRoute(router *gin.Engine, db *gorm.DB, validate *validator.Validate) {
	Products := service.NewProductService(
		repository.NewProductRepository(),
		db,
		validate,
	)
	productController := controller.NewProductController(Products)
	router.GET("/products", productController.FindAll)
	router.POST("/products", productController.Create)
}
