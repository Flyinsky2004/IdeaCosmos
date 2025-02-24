package service

import (
	"back-end/util"
	"context"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该更严格
	},
}

// type ChatRequest struct {
// 	Question  string `json:"question"`
// 	AuthToken string `json:"auth_token"`
// }

// func ChatStream(c *gin.Context) {
// 	upgrader := websocket.Upgrader{
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true
// 		},
// 	}
// 	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "WebSocket升级失败"))
// 		return
// 	}
// 	defer ws.Close()
// 	// 读取并解析初始消息
// 	_, message, err := ws.ReadMessage()
// 	if err != nil {
// 		return
// 	}
// 	var request ChatRequest
// 	if err := json.Unmarshal(message, &request); err != nil {
// 		ws.WriteJSON(dto.ErrorResponse[string](500, "无法解析请求参数"+err.Error()))
// 		return
// 	}

// 	// 验证token并获取userId
// 	claims, err := util.ParseToken(request.AuthToken)
// 	if err != nil {
// 		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
// 		return
// 	}
// 	userId := claims.UserID
// 	if userId == 0 {
// 		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
// 		return
// 	}
// 	// 调用流式聊天
// 	streamChan, err := util.StreamChatCompletion(ctx, util.ChatRequest{
// 		Model:       "deepseek-chat",
// 		Messages:    []util.Message{},
// 		Prompt:      systemPrompt,
// 		Question:    request.Question,
// 		Temperature: 0.6,
// 		MaxTokens:   8192,
// 	})

// 	if err != nil {
// 		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
// 		return
// 	}

// 	// 读取流式响应并通过WebSocket发送
// 	for response := range streamChan {
// 		if err := ws.WriteJSON(response); err != nil {
// 			return
// 		}
// 		if response.Done {
// 			break
// 		}
// 	}
// }

func HandleStreamChat(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// 读取初始请求
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}

	var chatRequest util.ChatRequest
	if err := json.Unmarshal(message, &chatRequest); err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "请求JSON有误" + err.Error(),
		})
		return
	}

	// 创建上下文，支持取消
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 获取流式响应
	streamChan, err := util.StreamChatCompletion(ctx, chatRequest)
	if err != nil {
		ws.WriteJSON(map[string]interface{}{
			"error": "流失请求失败" + err.Error(),
		})
		return
	}

	// 发送流式响应
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
	}
}
