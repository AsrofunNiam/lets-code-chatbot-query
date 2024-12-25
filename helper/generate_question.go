package helper

import (
	"context"
	"fmt"

	"github.com/AsrofunNiam/lets-code-chatbot-query/model/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateQuestion(apiKey, prompt string, contextChat *domain.ChatContext) (*genai.GenerateContentResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		PanicIfError(err)
	}
	defer client.Close()

	// add context
	if contextChat != nil {

		// entries prompt
		for _, entry := range contextChat.History {
			// prompt = fmt.Sprintf("%s\n\n%s: %s", prompt, entry.Type, entry.Content)
			prompt = fmt.Sprintf("%s\n\nMelanjutkan dari sebelumnya: %s", prompt, entry.Content)

		}

		//  all context
		// var additionalContext string
		// for _, entry := range contextChat.History {
		// 	additionalContext += fmt.Sprintf("\n\nMelanjutkan dari sebelumnya: %s", entry.Content)
		// }
		// prompt += additionalContext
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		PanicIfError(err)
	}
	return resp, nil
}
