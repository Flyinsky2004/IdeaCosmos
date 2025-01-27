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
