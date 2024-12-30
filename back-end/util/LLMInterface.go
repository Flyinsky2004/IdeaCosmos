package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	Temperature float64   `json:"temperature,omitempty"`
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

// Client 定义 ChatGPT API 客户端
type Client struct {
	apiKey  string
	baseURL string
}

// NewClient 创建一个新的 ChatGPT API 客户端
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: "https://api.deepseek.com/chat/completions",
	}
}

// SendMessage 发送消息并获取 AI 回复
func (c *Client) SendMessage(messages []Message, model string, maxToken int, temperature float64) (ChatResponse, error) {
	// 构建请求体
	reqBody := ChatRequest{
		Model:       model,
		Messages:    messages,
		MaxTokens:   maxToken,
		Temperature: temperature,
	}

	// 将请求体转换为 JSON
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to marshal request: %v", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to read response: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return ChatResponse{}, fmt.Errorf("API error: %s", body)
	}

	// 解析响应
	var chatResp ChatResponse
	err = json.Unmarshal(body, &chatResp)
	if err != nil {
		return ChatResponse{}, fmt.Errorf("failed to unmarshal response: %v", err)
	}
	return chatResp, nil
	// 返回第一个消息的内容
	//if len(chatResp.Choices) > 0 {
	//	return chatResp.Choices[0].Message.Content, nil
	//}

	return ChatResponse{}, fmt.Errorf("no response from API")
}

func ChatHandler(request ChatRequest) (ChatResponse, error) {
	client := NewClient("sk-46feaaa3624147ff9505aa8a8518dd6c")
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
