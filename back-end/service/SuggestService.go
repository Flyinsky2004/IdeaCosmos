/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 项目内容建议服务
 */
package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ProjectSuggestRequest struct {
	AuthToken string        `json:"auth_token"`
	ChatID    *uint         `json:"chat_id"`  // 可选，如果存在则表示继续之前的对话
	Messages  []ChatMessage `json:"messages"` // 历史消息记录
}

type ChatMessage struct {
	Role    string `json:"role"`    // 消息角色（system/user/assistant）
	Content string `json:"content"` // 消息内容
}

type ProjectInfo struct {
	ID          uint   `json:"id"`
	ProjectName string `json:"project_name"`
	Types       string `json:"types"`
	Style       string `json:"style"`
	SocialStory string `json:"social_story"`
	Watches     uint   `json:"watches"`
	Favorites   uint   `json:"favorites"`
	CreatedAt   string `json:"created_at"`
}

// ProjectContentSuggest 项目内容建议
func ProjectContentSuggest(c *gin.Context) {
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
	var request ProjectSuggestRequest
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
	userId := uint(claims.UserID)
	if userId == 0 {
		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		return
	}

	// 处理聊天会话
	var chat pojo.Chat
	if request.ChatID != nil {
		// 继续已有会话
		if err := config.MysqlDataBase.First(&chat, *request.ChatID).Error; err != nil {
			ws.WriteJSON(dto.ErrorResponse[string](500, "获取聊天会话失败"))
			return
		}
		// 验证会话所有权
		if chat.UserID != userId {
			ws.WriteJSON(dto.ErrorResponse[string](500, "无权访问该会话"))
			return
		}
	} else {
		// 创建新会话
		chat = pojo.Chat{
			UserID: userId,
			Type:   "project_suggest",
			Status: "active",
			Title:  "项目内容建议",
		}
		if err := config.MysqlDataBase.Create(&chat).Error; err != nil {
			ws.WriteJSON(dto.ErrorResponse[string](500, "创建聊天会话失败"))
			return
		}
	}

	// 获取最新的50条项目数据
	var projects []pojo.Project
	if err := config.MysqlDataBase.Model(&pojo.Project{}).
		Where("status = ?", "normal").
		Order("created_at DESC").
		Limit(50).
		Find(&projects).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取项目数据失败"+err.Error()))
		return
	}

	// 转换为简化的项目信息列表
	var projectInfos []ProjectInfo
	for _, project := range projects {
		var styleMap map[string]interface{}
		styleStr := "{}"
		if project.Style != nil {
			styleStr = string(project.Style)
		}
		json.Unmarshal([]byte(styleStr), &styleMap)

		styleDesc := ""
		if styleVal, ok := styleMap["name"]; ok {
			if str, ok := styleVal.(string); ok {
				styleDesc = str
			}
		}

		projectInfos = append(projectInfos, ProjectInfo{
			ID:          project.ID,
			ProjectName: project.ProjectName,
			Types:       project.Types,
			Style:       styleDesc,
			SocialStory: project.SocialStory,
			Watches:     project.Watches,
			Favorites:   project.Favorites,
			CreatedAt:   project.CreatedAt.Format("2006-01-02"),
		})
	}

	// 构建提示词
	prompt := `作为一个专业的影视内容分析师和创作顾问，请基于平台最新的50个项目数据，从创作者和观众双重视角进行分析。以下是平台上最新的项目数据（按创建时间排序）：
`
	for i, info := range projectInfos {
		prompt += fmt.Sprintf("%d. 项目名称：%s\n   类型：%s\n   风格：%s\n   社会故事：%s\n   观看次数：%d\n   收藏次数：%d\n   创建时间：%s\n\n",
			i+1, info.ProjectName, info.Types, info.Style, info.SocialStory, info.Watches, info.Favorites, info.CreatedAt)
	}

	prompt += `
请从以下两个视角进行分析：

观众视角：
1. 总结当前最受欢迎的作品特点（类型、风格、题材等）
2. 分析观众的观看和收藏偏好
3. 推荐3-5个最值得关注的优质项目，并说明推荐理由
4. 预测未来可能会受欢迎的内容方向

创作者视角：
1. 分析当前市场热点和创作趋势
2. 总结成功作品的共同特征（叙事手法、角色塑造等）
3. 提供3-5个具体的创作建议
4. 指出市场空白点和创新机会

请用专业但通俗易懂的语言回答，注重实用性，同时兼顾观众的观看体验和创作者的创作需求。回答中应该平衡艺术性与商业性，帮助创作者打造既有艺术价值又受欢迎的作品。`

	// 构建消息列表
	var messages []util.Message
	// 添加历史消息
	for _, msg := range request.Messages {
		messages = append(messages, util.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// 保存用户的新问题
	userMessage := pojo.ChatMessage{
		ChatID:  chat.ID,
		Role:    "user",
		Content: prompt,
		Status:  "success",
	}
	if err := config.MysqlDataBase.Create(&userMessage).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "保存用户消息失败"))
		return
	}

	// 调用流式聊天
	streamChan, err := util.StreamChatCompletion(context.Background(), util.ChatRequest{
		Model:       util.ThinkModelName,
		Messages:    messages,
		Prompt:      "你是一个专业的剧本创作顾问，擅长分析项目数据并提供专业的创作建议。",
		Question:    prompt,
		Temperature: util.GlobalTemperature,
		MaxTokens:   8192,
	})

	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 创建 AI 回复消息记录
	aiMessage := pojo.ChatMessage{
		ChatID:  chat.ID,
		Role:    "assistant",
		Content: "",
		Status:  "processing",
	}
	if err := config.MysqlDataBase.Create(&aiMessage).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "创建AI消息记录失败"))
		return
	}

	// 读取流式响应并通过WebSocket发送
	var fullResponse string
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
		if !response.Done {
			fullResponse += response.Content
		} else {
			// 更新AI消息内容和状态
			aiMessage.Content = fullResponse
			aiMessage.Status = "success"
			if err := config.MysqlDataBase.Save(&aiMessage).Error; err != nil {
				ws.WriteJSON(dto.ErrorResponse[string](500, "更新AI消息失败"))
				return
			}
			break
		}
	}
}
