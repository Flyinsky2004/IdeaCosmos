/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserChats 获取用户对话列表
func GetUserChats(c *gin.Context) {
	userIdOrigin, _ := c.Get("userId")
	userId := uint(userIdOrigin.(int))
	fmt.Println(userId)

	var chats []pojo.Chat
	if err := config.MysqlDataBase.Where("user_id = ?", userId).Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse[string](500, "获取对话列表失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(chats))
}

// GetChatHistory 获取对话历史记录
func GetChatHistory(c *gin.Context) {
	chatId := c.Param("chat_id")

	var messages []pojo.ChatMessage
	if err := config.MysqlDataBase.Where("chat_id = ?", chatId).Order("created_at ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse[string](500, "获取对话历史记录失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(messages))
}

// DeleteChat 删除对话记录
func DeleteChat(c *gin.Context) {
	chatId := c.Param("chat_id")

	if err := config.MysqlDataBase.Delete(&pojo.Chat{}, chatId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse[string](500, "删除对话记录失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("对话记录已删除"))
}
