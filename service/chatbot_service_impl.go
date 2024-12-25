package service

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
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
	ChatContexts     map[string]*domain.ChatContext
	mu               sync.Mutex
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
		ChatContexts:     make(map[string]*domain.ChatContext),
	}
}

func (service *ChatBotServiceImpl) Create(request *web.ChatBotCreateRequest, c *gin.Context) *genai.GenerateContentResponse {
	// tx := service.DB.Begin()
	// defer helper.CommitOrRollback(tx)

	tx := service.DB
	err := tx.Error
	helper.PanicIfError(err)

	// Load config
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}
	userID := c.GetString("userID")
	if userID == "" {
		userID = "default_user"
	}
	// Find context
	contextChat := service.GetChatContext(userID)

	if IsQueryRelated(request.Description) {
		descriptions := strings.Join(service.GenerateSchemaDescriptions(tx), "\n")
		prompt := fmt.Sprintf("Berdasarkan struktur tabel berikut:\n%s\nBuatkan query SQL untuk: %s. Hanya berikan query SQL tanpa penjelasan apapun atau format markdown.", descriptions, request.Description)

		resp, err := helper.GenerateQuestion(configuration.GeminiAPIKey, prompt, contextChat)
		if err != nil {
			log.Fatalf("Failed to generate content: %v", err)
		}

		query := helper.ExtractQuery(resp.Candidates[0].Content.Parts[0])
		var resultValueDb []map[string]interface{}
		err = tx.Raw(query).Find(&resultValueDb).Error
		if err != nil {
			helper.PanicIfError(err)
		}

		resultString := formatQueryResult(resultValueDb)
		newPrompt := fmt.Sprintf("Hasil query adalah: %s. Jelaskan hasil ini dengan bahasa manusia.", resultString)

		respDesc, err := helper.GenerateQuestion(configuration.GeminiAPIKey, newPrompt, contextChat)
		if err != nil {
			log.Fatalf("Failed to generate description: %v", err)
		}
		contextChat.History = append(contextChat.History, domain.ChatHistoryEntry{
			Type:      "query",
			Content:   newPrompt,
			Timestamp: time.Now(),
		})

		fmt.Println("memory 1 ", contextChat.History)
		fmt.Println("len memory 2 ", len(contextChat.History))

		return respDesc
	} else {
		// Prompt to generate response regular
		newPrompt := fmt.Sprintf("Jawab pertanyaan ini berdasarkan konteks manusia: %s", request.Description)
		resp, err := helper.GenerateQuestion(configuration.GeminiAPIKey, newPrompt, contextChat)
		if err != nil {
			log.Fatalf("Failed to generate response: %v", err)
		}

		// Add the new prompt to the history
		contextChat.History = append(contextChat.History, domain.ChatHistoryEntry{
			Type:      "general",
			Content:   newPrompt,
			Timestamp: time.Now(),
		})

		fmt.Println("memory general ", contextChat.History)
		fmt.Println("len memory general ", len(contextChat.History))
		return resp
	}
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

func IsQueryRelated(question string) bool {
	keywords := []string{"query", "SQL", "database", "tabel"}
	for _, keyword := range keywords {
		if strings.Contains(strings.ToLower(question), keyword) {
			return true
		}
	}
	return false
}

func (service *ChatBotServiceImpl) GetChatContext(userID string) *domain.ChatContext {
	service.mu.Lock()
	defer service.mu.Unlock()

	// Check if context already exists
	if context, exists := service.ChatContexts[userID]; exists {
		return context
	}

	// If context doesn't exist, create a new one
	newContext := &domain.ChatContext{
		UserID:  userID,
		History: []domain.ChatHistoryEntry{},
	}
	service.ChatContexts[userID] = newContext
	return newContext
}
