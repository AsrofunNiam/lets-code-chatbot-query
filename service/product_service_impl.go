package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/web"
	"github.com/AsrofunNiam/lets-code-chatbot-query/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/generative-ai-go/genai"
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

func (service *ProductServiceImpl) Create(request *web.ProductCreateRequest, c *gin.Context) *genai.GenerateContentResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Load config
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	// Generate query schema database descriptions
	descriptions := strings.Join(helper.GenerateSchemaDescriptions(tx), "\n")

	//  Generate prompt first
	prompt := fmt.Sprintf("Berdasarkan struktur tabel berikut:\n%s\nBuatkan query SQL untuk: %s. Hanya berikan query SQL tanpa penjelasan apapun atau format markdown.", descriptions, request.Description)
	resp, err := helper.GenerateQuestion(configuration.GeminiAPIKey, prompt)
	if err != nil {
		log.Fatalf("Failed to generate content: %v", err)
	}

	// convert interface to string
	query := helper.ExtractQuery(resp.Candidates[0].Content.Parts[0])

	// Execute query
	var resultValueDb []map[string]interface{}
	err = tx.Raw(query).Scan(&resultValueDb).Error
	if err != nil {
		helper.PanicIfError(err)
	}

	// Send prompt second
	resultString := formatQueryResult(resultValueDb)
	newPrompt := fmt.Sprintf("Hasil query adalah: %s. Jelaskan hasil ini dengan bahasa manusia.", resultString)

	// Generate description second
	respDesc, err := helper.GenerateQuestion(configuration.GeminiAPIKey, newPrompt)
	if err != nil {
		log.Fatalf("Failed to generate description: %v", err)
	}

	return respDesc
}

// Helper function untuk format hasil query menjadi string
func formatQueryResult(results []map[string]interface{}) string {
	var sb strings.Builder
	for _, row := range results {
		for key, value := range row {
			sb.WriteString(fmt.Sprintf("%s: %v, ", key, value))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
