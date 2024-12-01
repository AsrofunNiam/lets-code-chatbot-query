package helper

import (
	"context"
	"fmt"

	// "os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func GenerateQuestion(apiKey, prompt string) (*genai.GenerateContentResponse, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		PanicIfError(err)
	}
	defer client.Close()

	// [START text_gen_text_only_prompt]
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		PanicIfError(err)
	}

	return resp, nil
}

func PrintResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
