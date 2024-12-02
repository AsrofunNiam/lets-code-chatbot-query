package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AsrofunNiam/lets-code-chatbot-query/app"
	c "github.com/AsrofunNiam/lets-code-chatbot-query/configuration"
	"github.com/AsrofunNiam/lets-code-chatbot-query/helper"
	"github.com/go-playground/validator/v10"
)

func main() {
	// ctx := context.Background()
	configuration, err := c.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	//  connection db
	port := configuration.Port
	db := app.ConnectDatabase(configuration.User, configuration.Host, configuration.Password, configuration.PortDB, configuration.Db)

	// connection db
	// client, err := genai.NewClient(context.Background(), genai.WithAPIKey(configuration.GeminiAPIKey))
	// if err != nil {
	// 	log.Fatalf("Failed to create Gemini client: %v", err)
	// }

	// client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv(configuration.GeminiAPIKey)))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// geminiClient = app.SetGeminiClient(client)

	fmt.Print("API Key: ", configuration.GeminiAPIKey)

	validate := validator.New()
	router := app.NewRouter(db, validate)
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("Server is running on port %s", port)

	err = server.ListenAndServe()
	helper.PanicIfError(err)
}
