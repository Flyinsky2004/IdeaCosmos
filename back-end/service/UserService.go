/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetMyInfo(c *gin.Context) {
	var user pojo.User
	userId, _ := c.Get("userId")
	if err := config.MysqlDataBase.Where("id = ?", userId).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在获取用户信息时出错！详细信息:"+err.Error()))
		return
	}
	user.Password = ""
	c.JSON(http.StatusOK, dto.SuccessResponse(user))
}

// GetAllUsers 获取所有用户列表（仅ID、用户名、头像）
func GetAllUsers(c *gin.Context) {
	type SimpleUser struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	var users []SimpleUser

	// 查询所有用户，只获取ID、用户名和头像字段
	if err := config.MysqlDataBase.Model(&pojo.User{}).
		Select("id, username, avatar").
		Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取用户列表失败:"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"users": users,
	}))
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,max=200"`
}

func UpdateUserInfo(c *gin.Context) {
	userID, _ := c.Get("userId")
	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在获取用户信息时出错！详细信息:"+err.Error()))
		return
	}

	// 解析请求体
	var updateRequest UpdateUserRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 更新用户信息
	if updateRequest.Username != "" {
		user.Username = updateRequest.Username
	}
	if updateRequest.Avatar != "" {
		user.Avatar = updateRequest.Avatar
	}
	tx := config.MysqlDataBase.Begin()
	// 保存到数据库
	if err := tx.Save(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在保存用户信息时出错！详细信息:"+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在保存用户信息时出错！详细信息:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "用户信息更新成功！"))
}
