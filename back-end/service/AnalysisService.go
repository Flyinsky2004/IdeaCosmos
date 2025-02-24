package service

import (
	"back-end/entity/dto"
	"back-end/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NewProjectAnalysisRequest struct {
	Data      []StyleStats `json:"data"`
	AuthToken string       `json:"auth_token"`
}

// NewProjectAnalysis 项目分析建议
func NewProjectAnalysis(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "WebSocket升级失败"))
		return
	}
	defer ws.Close()

	// 读取并解析初始消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}
	var request NewProjectAnalysisRequest
	if err := json.Unmarshal(message, &request); err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "无法解析请求参数"+err.Error()))
		return
	}

	// 验证token并获取userId
	claims, err := util.ParseToken(request.AuthToken)
	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		return
	}
	userId := claims.UserID
	if userId == 0 {
		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		return
	}

	// 按观看数和收藏数排序并获取前三
	type StyleAnalysis struct {
		Style      string
		Type       string
		WatchCount int
		LikeCount  int
		Total      int // 用于排序的总分
	}

	// 聚合数据
	styleMap := make(map[string]*StyleAnalysis)
	for _, stat := range request.Data {
		key := stat.Style + "|" + stat.Type
		if _, exists := styleMap[key]; !exists {
			styleMap[key] = &StyleAnalysis{
				Style: stat.Style,
				Type:  stat.Type,
			}
		}
		styleMap[key].WatchCount += stat.WatchCount
		styleMap[key].LikeCount += stat.LikeCount
		// 计算总分 (这里可以调整权重)
		styleMap[key].Total = styleMap[key].WatchCount + styleMap[key].LikeCount*2 // 收藏权重更高
	}

	// 转换为切片并排序
	var analyses []StyleAnalysis
	for _, analysis := range styleMap {
		analyses = append(analyses, *analysis)
	}
	sort.Slice(analyses, func(i, j int) bool {
		return analyses[i].Total > analyses[j].Total
	})

	// 获取前三名
	topThree := analyses
	if len(topThree) > 3 {
		topThree = topThree[:3]
	}

	// 构建提示词
	prompt := `作为一个专业的剧本分析师，我需要你分析以下数据并给出建议。

根据最近的数据统计，用户最喜欢的前三个风格和类型组合是：

`
	for i, analysis := range topThree {
		prompt += fmt.Sprintf("%d. 风格：%s，类型：%s\n   观看次数：%d，收藏次数：%d\n\n",
			i+1, analysis.Style, analysis.Type, analysis.WatchCount, analysis.LikeCount)
	}

	prompt += `
请你：
1. 简要总结这些数据反映出的用户偏好
2. 基于这些最受欢迎的风格和类型组合，提出3-5个具体的剧集创作建议
3. 分析这些风格和类型组合为什么会受欢迎，以及如何在创作中更好地运用它们

请用专业但易懂的语言回答，并注重实用性建议。`

	// 调用流式聊天
	streamChan, err := util.StreamChatCompletion(context.Background(), util.ChatRequest{
		Model:       "deepseek-chat",
		Messages:    []util.Message{},
		Prompt:      "你是一个专业的剧本分析师，擅长分析数据并提供专业的创作建议。",
		Question:    prompt,
		Temperature: 0.7,
		MaxTokens:   8096,
	})

	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 读取流式响应并通过WebSocket发送
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
		if response.Done {
			break
		}
	}
}
