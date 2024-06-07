# perplexity-go

perplexity-go is a non-official Go client library for the [Perplexity AI](https://www.perplexity.ai/) chat completion [API](https://docs.perplexity.ai/reference/post_chat_completions). It provides an easy-to-use interface for interacting with the Perplexity API and generating chat completions using various language models.

> [!WARNING]
> Experimental package, API would probably break.

## Features

- Simple and intuitive API for chat completions
- Support for multiple language models
- Customizable request parameters
- Asynchronous requests with context cancellation

## Not supported yet

- Streaming
- Images
- Citations

## Installation

To install perplexity-go, use `go get`:

```bash
go get -u github.com/painhardcore/perplexity-go
```

## Example Usage
```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/painhardcore/perplexity-go"
)

func main() {
	apiKey := "pplx-xxx" // Replace with your actual API key
	model := perplexity.ModelLlama3SonarSmall32kChat

	client := perplexity.NewClient(apiKey, model)

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
```

## Contributing

Contributions to perplexity-go are welcome! If you find a bug or have a feature request, please open an issue on the GitHub repository. If you'd like to contribute code, please fork the repository and submit a pull request.
