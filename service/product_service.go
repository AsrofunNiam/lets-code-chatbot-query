package service

import (
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/gin-gonic/gin"
)

type ProductService interface {
	FindAll(filters *map[string]string, c *gin.Context) []web.ProductResponse
	Create(request *web.ProductCreateRequest, c *gin.Context) web.ProductResponse
}
