package util

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/29 17:25
 */
// ChatRequest 定义发送给 ChatGPT API 的请求结构
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Prompt      string    `json:"prompt,omitempty"`
	Question    string    `json:"question,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
}

// Message 定义消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse 定义 ChatGPT API 返回的响应结构
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Client 使用 OpenAI SDK 客户端
type Client struct {
	client  *openai.Client
	baseURL string
}

// NewClient 创建一个新的客户端
func NewClient(apiKey string) *Client {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = OpenAIBaseURL // 使用与 stream_chat.go 相同的 URL

	return &Client{
		client:  openai.NewClientWithConfig(config),
		baseURL: OpenAIBaseURL,
	}
}

// SendMessage 使用 OpenAI SDK 发送消息
func (c *Client) SendMessage(messages []Message, model string, maxToken int, temperature float32) (ChatResponse, error) {
	// 转换消息格式
	openaiMessages := make([]openai.ChatCompletionMessage, len(messages))
	for i, msg := range messages {
		openaiMessages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 创建请求
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       UseModelName,
			Messages:    openaiMessages,
			MaxTokens:   maxToken,
			Temperature: float32(temperature),
		},
	)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to create chat completion: %v", err)
	}

	// 转换响应格式为原有的 ChatResponse 结构
	chatResp := ChatResponse{
		Choices: []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		}{
			{
				Message: struct {
					Content string `json:"content"`
				}{
					Content: resp.Choices[0].Message.Content,
				},
			},
		},
		Usage: struct {
			TotalTokens int `json:"total_tokens"`
		}{
			TotalTokens: resp.Usage.TotalTokens,
		},
	}

	return chatResp, nil
}

// ChatHandler 保持原有函数签名和行为不变
func ChatHandler(request ChatRequest) (ChatResponse, error) {
	client := NewClient(OpenAIKey) // 使用与 stream_chat.go 相同的 key

	systemMessage := []Message{
		{
			Role:    "system",
			Content: request.Prompt,
		},
	}
	askMessage := []Message{
		{
			Role:    "user",
			Content: request.Question,
		},
	}
	messages := append(systemMessage, request.Messages...)
	messages = append(messages, askMessage...)

	response, err := client.SendMessage(messages, request.Model, request.MaxTokens, request.Temperature)
	if err != nil {
		fmt.Println("Error:", err)
		return ChatResponse{}, err
	}

	return response, nil
}
