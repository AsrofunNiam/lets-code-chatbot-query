package controller

import (
	"net/http"

	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/AsrofunNiam/lets-code-chatbot-query/service"
	"github.com/gin-gonic/gin"
)

type ProductControllerImpl struct {
	ProductService service.ProductService
}

func NewProductController(productService service.ProductService) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
	}
}

func (controller *ProductControllerImpl) FindAll(c *gin.Context) {
	filters := helper.FilterFromQueryString(c, "name.like", "id.eq")
	productResponses := controller.ProductService.FindAll(&filters, c)
	webResponse := web.WebResponse{
		Success: true,
		Message: helper.MessageDataFoundOrNot(productResponses),
		Data:    productResponses,
	}

	c.JSON(http.StatusOK, webResponse)
}

func (controller *ProductControllerImpl) Create(c *gin.Context) {
	request := &web.ProductCreateRequest{}
	helper.ReadFromRequestBody(c, &request)

	productResponse := controller.ProductService.Create(request, c)
	webResponse := web.WebResponse{
		Success: true,
		Message: "Product created successfully",
		Data:    productResponse,
	}

	c.JSON(http.StatusOK, webResponse)
}
