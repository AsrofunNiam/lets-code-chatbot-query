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

type ChatBotServiceImpl struct {
	SchemaRepository repository.SchemaRepository
	DB               *gorm.DB
	Validate         *validator.Validate
}

func NewChatBotService(
	schemaRepository repository.SchemaRepository,
	db *gorm.DB,
	validate *validator.Validate,
) ChatBotService {
	return &ChatBotServiceImpl{
		SchemaRepository: schemaRepository,
		DB:               db,
		Validate:         validate,
	}
}

func (service *ChatBotServiceImpl) Create(request *web.ChatBotCreateRequest, c *gin.Context) *genai.GenerateContentResponse {
	tx := service.DB.Begin()
	defer helper.CommitOrRollback(tx)

	// Load config
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	// Generate query schema database descriptions
	descriptions := strings.Join(service.GenerateSchemaDescriptions(tx), "\n")

	//  Generate prompt first
	prompt := fmt.Sprintf("Based on the following table structure:\n%s\nGenerate an SQL query for: %s. Only provide the SQL query without any explanation or markdown format.", descriptions, request.Description)
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
	newPrompt := fmt.Sprintf("The result of the query is: %s. Explain this result in provide the explanation in Indonesia and human-readable language.", resultString)

	// Generate description second
	respDesc, err := helper.GenerateQuestion(configuration.GeminiAPIKey, newPrompt)
	if err != nil {
		log.Fatalf("Failed to generate description: %v", err)
	}

	return respDesc
}

func (service *ChatBotServiceImpl) GenerateSchemaDescriptions(tx *gorm.DB) []string {
	databaseSchema := service.SchemaRepository.FindAll(tx)
	var descriptions []string
	for _, col := range databaseSchema {
		desc := fmt.Sprintf("Table: %s, Column: %s, Type: %s", col.TableName, col.ColumnName, col.DataType)
		descriptions = append(descriptions, desc)
	}

	return descriptions
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
