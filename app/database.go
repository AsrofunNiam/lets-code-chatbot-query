package app

import (
	"log"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabase(user, host, password, port, db string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + db + "?parseTime=true"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}

	//  function auto migrate, generate schema table
	err = database.AutoMigrate(
	// &domain.User{},
	)
	if err != nil {
		panic("failed to auto migrate schema")
	}

	return database
}

var geminiClient *genai.Client

// SetGeminiClient sets the global Gemini client
func SetGeminiClient(client *genai.Client) *genai.Client {
	geminiClient = client

	return geminiClient
}

// GetGeminiClient returns the global Gemini client
func GetGeminiClient() *genai.Client {
	if geminiClient == nil {
		log.Fatal("Gemini client is not initialized")
	}
	return geminiClient
}
