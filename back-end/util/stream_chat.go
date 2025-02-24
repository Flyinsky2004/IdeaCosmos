package util

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

const (
	UseModelName = "deepseek-r1-250120"
	// // deepseek官方
	// OpenAIKey = "sk-dd6f248183a846d692fe78b075696676"
	// OpenAIBaseURL = "https://api.deepseek.com/v1"
	// //1024110中转
	// OpenAIKey     = "sk-XnbHbzBOmPYGHgL_UV9xTM47qLbJs44jbIjK12f31ggQfX1MRFHCbRNwOec"
	// OpenAIBaseURL = "https://models.1024110.xyz/v1"
	//火山大模型中转
	OpenAIKey     = "817b8b32-9ccc-4055-aab0-5f4fe97025d6"
	OpenAIBaseURL = "https://ark.cn-beijing.volces.com/api/v3"
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
			Model:       UseModelName,
			Messages:    messages,
			Temperature: request.Temperature,
			Stream:      true,
			MaxTokens:   request.MaxTokens,
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
