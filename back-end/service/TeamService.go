package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/28 16:48
 */
func CreateTeam(c *gin.Context) {
	var team pojo.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}
	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在创建团队信息时出错！详细信息:"+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "在存储团队信息时出错！详细信息:"+err.Error()))
		tx.Rollback()
		return
	}
	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "团队创建成功！"))
}
func UpdateTeam(c *gin.Context) {
	userId, _ := c.Get("userId")
	var reqBody pojo.TeamUpdateBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
	}
	tx := config.MysqlDataBase.Begin()
	var team pojo.Team
	if err := tx.First(&team, reqBody.ID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "找不到请求的团队"))
		return
	}

	if team.LeaderId != userId {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改该团队的信息"))
		return
	}

	team.TeamName = reqBody.TeamName
	team.TeamDescription = reqBody.TeamDescription

	if err := tx.Save(&team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存团队信息时发生错误："+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存团队信息时发生错误："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "团队信息更新成功！"))
}
func GetMyTeam(c *gin.Context) {
	userId, _ := c.Get("userId")
	offset := c.Query("offset")
	offsetInt, _ := strconv.Atoi(offset)
	var teams []pojo.Team
	if err := config.MysqlDataBase.Where("leader_id = ?", userId).Limit(10).Offset(offsetInt).Find(&teams).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询团队信息时发生错误："+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[[]pojo.Team](teams))
}
func GetMyJoinedTeam(c *gin.Context) {
	var requests []pojo.JoinRequest
	offset := c.PostForm("offset")
	offsetInt, _ := strconv.Atoi(offset)
	userId, _ := c.Get("userId")
	if err := config.MysqlDataBase.Preload("Team").
		Joins("JOIN teams ON join_requests.team_id = teams.id").
		Where("teams.leader_id = ? AND status = ?", userId, 1).
		Offset(offsetInt).
		Limit(10).
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队信息时发生错误！"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[[]pojo.JoinRequest](requests))
}

const (
	StatusPending  = 0 // Pending
	StatusApproved = 1 // Approved
	StatusRejected = 2 // Rejected
)

func RequestToJoin(c *gin.Context) {
	var joinRequest pojo.JoinRequest
	if err := c.ShouldBindJSON(&joinRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交表单格式有误！"))
		return
	}

	// 校验团队是否存在
	var team pojo.Team
	tx := config.MysqlDataBase.Begin()
	if err := tx.First(&team, joinRequest.TeamId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "请求团队不存在！"))
		return
	}

	// 默认状态为 pending (0)
	joinRequest.Status = StatusPending
	if err := tx.Create(&joinRequest).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存团队加入申请时发生错误！"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存团队加入申请时发生错误！"))
		return
	}

	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "团队加入申请发送成功！"))
}

func UpdateRequest(c *gin.Context) {
	userId, _ := c.Get("userId")
	var reqBody pojo.UpdateJoinRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交表单格式有误！"))
		return
	}
	var joinRequest pojo.JoinRequest
	tx := config.MysqlDataBase.Begin()
	if err := tx.First(&joinRequest, reqBody.RequestId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "请求团队请求不存在！"))
		return
	}
	var team pojo.Team
	if err := tx.First(&team, joinRequest.TeamId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "请求团队信息不存在！"))
		return
	}
	if team.LeaderId != userId {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改该团队的信息"))
		return
	}
	joinRequest.Status = reqBody.Status
	if err := tx.Save(&joinRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存指定团队申请时发生错误！"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新指定团队申请时发生错误！"))
		return
	}
	c.JSON(http.StatusOK, dto.ErrorResponse[string](200, "团队请求更新成功！"))
}
func GetPendingRequests(c *gin.Context) {
	var requests []pojo.JoinRequest
	offset := c.PostForm("offset")
	status := c.PostForm("status")
	offsetInt, _ := strconv.Atoi(offset)
	userId, _ := c.Get("userId")
	if err := config.MysqlDataBase.Preload("Team").
		Joins("JOIN teams ON join_requests.team_id = teams.id").
		Where("teams.leader_id = ? AND status = ?", userId, status).
		Offset(offsetInt).
		Limit(10).
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队申请时发生错误！"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[[]pojo.JoinRequest](requests))
}
