package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 14:40
 */

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

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,max=50"`
	Avatar   string `json:"avatar" binding:"omitempty,max=80"`
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
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在保存用户信息时出错！详细信息:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "用户信息更新成功！"))
}
