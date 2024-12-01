package service

import (
	"fmt"
	"log"

	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/AsrofunNiam/lets-code-chatbot-query/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	config "github.com/AsrofunNiam/lets-code-chatbot-query/configuration"
)

type ProductServiceImpl struct {
	ProductRepository repository.ProductRepository
	DB                *gorm.DB
	Validate          *validator.Validate
}

func NewProductService(
	product repository.ProductRepository,
	db *gorm.DB,
	validate *validator.Validate,
) ProductService {
	return &ProductServiceImpl{
		ProductRepository: product,
		DB:                db,
		Validate:          validate,
	}
}

func (service *ProductServiceImpl) FindAll(filters *map[string]string, c *gin.Context) []web.ProductResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	products := service.ProductRepository.FindAll(tx, filters)
	return products.ToProductResponses()
}

func (service *ProductServiceImpl) Create(request *web.ProductCreateRequest, c *gin.Context) web.ProductResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	filters := map[string]string{"name.eq": request.Name}

	//  resource question
	products := service.ProductRepository.FindAll(tx, &filters)

	fmt.Println(len(products))

	resp, err := helper.GenerateQuestion(configuration.GeminiAPIKey, request.Description)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	// Optionally, you can use the response here or pass it to another service
	helper.PrintResponse(resp)

	product := &domain.Product{
		// Required Fields
		Name:        request.Name,
		Description: request.Description,
	}
	product = service.ProductRepository.Create(tx, product)
	return product.ToProductResponse()
}
