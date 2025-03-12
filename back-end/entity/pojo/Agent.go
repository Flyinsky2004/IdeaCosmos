/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: 对话历史记录相关模型
 */
package pojo

import (
	"gorm.io/gorm"
)

// Chat 聊天会话模型
type Chat struct {
	gorm.Model
	UserID   uint          `gorm:"not null;index" json:"user_id"`     // 用户ID
	Type     string        `gorm:"type:varchar(50)" json:"type"`      // 对话类型（如：project_suggest）
	Status   string        `gorm:"type:varchar(20)" json:"status"`    // 会话状态（如：active, completed, failed）
	Title    string        `gorm:"type:varchar(255)" json:"title"`    // 会话标题
	Messages []ChatMessage `gorm:"foreignKey:ChatID" json:"messages"` // 关联的消息列表
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	gorm.Model
	ChatID     uint   `gorm:"not null;index" json:"chat_id"`  // 关联的聊天会话ID
	Role       string `gorm:"type:varchar(20)" json:"role"`   // 消息角色（system/user/assistant）
	Content    string `gorm:"type:text" json:"content"`       // 消息内容
	TokenCount int    `gorm:"default:0" json:"token_count"`   // token数量
	Status     string `gorm:"type:varchar(20)" json:"status"` // 消息状态（如：success, failed）
}
