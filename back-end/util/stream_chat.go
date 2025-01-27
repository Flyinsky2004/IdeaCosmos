package util

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

const (
	// OpenAIKey OpenAI API密钥
	OpenAIKey = "sk-46feaaa3624147ff9505aa8a8518dd6c"
	// OpenAIBaseURL OpenAI API基础URL
	OpenAIBaseURL = "https://api.deepseek.com/v1"
)

type StreamResponse struct {
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func StreamChatCompletion(ctx context.Context, request ChatRequest) (<-chan StreamResponse, error) {
	// 创建配置
	config := openai.DefaultConfig(OpenAIKey)
	config.BaseURL = OpenAIBaseURL

	// 使用配置创建客户端
	client := openai.NewClientWithConfig(config)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    "system",
			Content: request.Prompt,
		},
		{
			Role:    "user",
			Content: request.Question,
		},
	}

	stream, err := client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:       request.Model,
			Messages:    messages,
			Temperature: request.Temperature,
			MaxTokens:   request.MaxTokens,
			Stream:      true,
		},
	)
	if err != nil {
		return nil, err
	}

	responseChan := make(chan StreamResponse)

	go func() {
		defer stream.Close()
		defer close(responseChan)

		for {
			response, err := stream.Recv()
			if err != nil {
				// 流结束
				responseChan <- StreamResponse{
					Done: true,
				}
				return
			}

			if len(response.Choices) > 0 {
				responseChan <- StreamResponse{
					Content: response.Choices[0].Delta.Content,
					Done:    false,
				}
			}
		}
	}()

	return responseChan, nil
}
