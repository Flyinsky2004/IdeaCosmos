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

// IdeaCosmosChat 创剧星球原住民
func IdeaCosmosChat(c *gin.Context) {
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

		// 加载历史消息
		if err := config.MysqlDataBase.Model(&chat).Association("Messages").Find(&chat.Messages); err != nil {
			ws.WriteJSON(dto.ErrorResponse[string](500, "获取历史消息失败"))
			return
		}
	} else {
		// 创建新会话
		chat = pojo.Chat{
			UserID: userId,
			Type:   "idea_cosmos_chat",
			Status: "active",
			Title:  "创剧星球助手对话",
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
	prompt := `我是创剧星球的智能助手，作为这个专业的剧集内容创作平台的原住民，我可以：

- 为创作者提供创意启发和市场洞察
- 为观众推荐优质内容和个性化建议
- 解答平台使用相关的各类问题
- 提供行业动态和趋势分析

以下是平台上最新的优质项目数据（按创建时间排序）：
`
	for i, info := range projectInfos {
		prompt += fmt.Sprintf("%d. 项目名称：%s\n   类型：%s\n   风格：%s\n   社会故事：%s\n   观看次数：%d\n   收藏次数：%d\n   创建时间：%s\n\n",
			i+1, info.ProjectName, info.Types, info.Style, info.SocialStory, info.Watches, info.Favorites, info.CreatedAt)
	}

	prompt += `
基于这些数据和我的专业知识，我可以为您提供以下服务：

创作者服务：
1. 剧本创意与故事构建建议
2. 市场趋势与受众分析
3. 项目优化和差异化建议
4. 创作技巧与经验分享

观众服务：
1. 个性化内容推荐
2. 热门作品解析
3. 类似作品推荐
4. 观看建议与导航

平台服务：
1. 功能介绍与使用指南
2. 创作工具使用建议
3. 社区互动与反馈
4. 平台活动与资源推荐

行业洞察：
1. 剧集市场趋势分析
2. 用户行为与偏好研究
3. 新兴题材与机会点
4. 行业动态与发展方向

请告诉我您的需求，我会以专业、友好的方式为您提供帮助。无论您是创作者还是观众，或者对平台有任何疑问，我都很乐意为您服务。`

	// 构建消息列表
	var messages []util.Message

	// 如果是继续对话，添加历史消息到messages
	if request.ChatID != nil {
		for _, msg := range chat.Messages {
			messages = append(messages, util.Message{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// 添加新的消息
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
		Content: request.Messages[len(request.Messages)-1].Content, // 使用用户的实际问题
		Status:  "processing",                                      // 初始状态为处理中
	}
	if err := config.MysqlDataBase.Create(&userMessage).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "保存用户消息失败"))
		return
	}

	// 调用流式聊天
	streamChan, err := util.StreamChatCompletion(context.Background(), util.ChatRequest{
		Model:       util.ThinkModelName,
		Messages:    messages,
		Prompt:      prompt,
		Question:    request.Messages[len(request.Messages)-1].Content, // 使用用户的实际问题
		Temperature: util.GlobalTemperature,
		MaxTokens:   8192,
	})

	if err != nil {
		// 更新用户消息状态为失败
		userMessage.Status = "failed"
		config.MysqlDataBase.Save(&userMessage)
		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 更新用户消息状态为成功
	userMessage.Status = "success"
	if err := config.MysqlDataBase.Save(&userMessage).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "更新用户消息状态失败"))
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
			// 更新AI消息状态为失败
			aiMessage.Status = "failed"
			config.MysqlDataBase.Save(&aiMessage)
			return
		}
		if !response.Done {
			fullResponse += response.Content
		} else {
			// 更新AI消息内容和状态
			aiMessage.Content = fullResponse
			aiMessage.Status = "success"

			// 如果是新对话，更新对话标题
			if request.ChatID == nil {
				// 从用户问题和AI回答中提取标题
				title := generateChatTitle(request.Messages[len(request.Messages)-1].Content, fullResponse)
				chat.Title = title
			}

			if err := config.MysqlDataBase.Save(&aiMessage).Error; err != nil {
				ws.WriteJSON(dto.ErrorResponse[string](500, "更新AI消息失败"))
				return
			}

			// 更新对话状态和标题
			chat.Status = "active"
			if err := config.MysqlDataBase.Save(&chat).Error; err != nil {
				ws.WriteJSON(dto.ErrorResponse[string](500, "更新对话状态失败"))
				return
			}
			break
		}
	}
}

// generateChatTitle 生成对话标题
func generateChatTitle(question, answer string) string {
	// 获取问题的前15个字符
	questionPart := ""
	if len([]rune(question)) > 15 {
		questionPart = string([]rune(question)[:15])
	} else {
		questionPart = question
	}

	// 获取答案的前15个字符
	answerPart := ""
	if len([]rune(answer)) > 15 {
		answerPart = string([]rune(answer)[:15])
	} else {
		answerPart = answer
	}

	// 组合标题，确保总长度不超过30个字符
	title := questionPart
	remainingLength := 30 - len([]rune(questionPart))
	if remainingLength > 0 && len(answerPart) > 0 {
		title += "..."
		if remainingLength > 3 {
			if len([]rune(answerPart)) > remainingLength-3 {
				title += string([]rune(answerPart)[:remainingLength-3])
			} else {
				title += answerPart
			}
		}
	}

	return title
}
