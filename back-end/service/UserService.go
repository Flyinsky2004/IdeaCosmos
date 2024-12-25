package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"github.com/gin-gonic/gin"
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
		c.JSON(500, dto.ErrorResponse[string](500, "在获取用户信息时出错！详细信息:"+err.Error()))
		return
	}
	user.Password = ""
	c.JSON(200, dto.SuccessResponse(user))
}
