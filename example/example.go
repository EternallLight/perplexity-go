package main

import (
	"context"
	"fmt"
	"log"

	"github.com/painhardcore/perplexity-go"
)

func main() {
	apiKey := "pplx-xxx" // Replace with your actual API key
	model := perplexity.ModelLlama31SonarSmall128kChat

	client := perplexity.NewClient(apiKey, model, nil)

	request := perplexity.ChatCompletionRequest{
		Messages: []perplexity.Message{
			{Role: perplexity.RoleSystem, Content: "Be rude to me, but respect me"},
			{Role: perplexity.RoleUser, Content: "Tell me how to be a good dev and solve all tickets in my sprints"},
		},
		MaxTokens:   100,
		Temperature: 0.2,
		TopP:        0.9,
	}

	ctx := context.Background()
	response, err := client.ChatCompletions(ctx, request)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	text, err := response.GetCompleteSingleMessage()
	if err != nil {
		fmt.Printf("err: %s,\nfull response %+v\n", err, response)
	}
	fmt.Println(text)
}
