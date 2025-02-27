package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/28 16:48
 */
type CreateTeamRequestBody struct {
	gorm.Model
	TeamName        string `json:"team_name" gorm:"type:varchar(50)"`
	TeamDescription string `json:"team_description" gorm:"type:varchar(50)"`
}
type TeamWithCount struct {
	pojo.Team
	MemberCount int `json:"member_count"`
}

func teamBodyConvertor(teams []pojo.Team) []TeamWithCount {
	var teamsWithCount []TeamWithCount
	for _, team := range teams {
		var count int64
		config.MysqlDataBase.Model(&pojo.JoinRequest{}).
			Where("team_id = ? AND status = ?", team.ID, StatusApproved).
			Count(&count)

		teamsWithCount = append(teamsWithCount, TeamWithCount{
			Team:        team,
			MemberCount: int(count) + 1, // +1 for the leader
		})
	}
	return teamsWithCount
}
func GetUserTeams(c *gin.Context) {
	userId, _ := c.Get("userId")
	var teams []pojo.Team

	// 子查询：用户创建的团队
	createdTeamsSubQuery := config.MysqlDataBase.Model(&pojo.Team{}).
		Select("id").
		Where("leader_id = ?", userId)

	// 子查询：用户加入的团队
	joinedTeamsSubQuery := config.MysqlDataBase.Model(&pojo.JoinRequest{}).
		Select("team_id").
		Where("user_id = ? AND status = ?", userId, 1)

	// 查询用户相关的团队
	err := config.MysqlDataBase.Where("id IN (?) OR id IN (?)", createdTeamsSubQuery, joinedTeamsSubQuery).
		Find(&teams).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "数据库查询发生错误"))
		return
	}

	res := teamBodyConvertor(teams)
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

func CreateTeam(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)
	var reqBody CreateTeamRequestBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}
	team := pojo.Team{
		TeamName:        reqBody.TeamName,
		TeamDescription: reqBody.TeamDescription,
		LeaderId:        uint(userIdInt),
		InviteCode:      util.GenerateRandomString(8),
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
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("团队创建成功！", "团队创建成功！"))
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
	res := teamBodyConvertor(teams)
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}
func GetMyJoinedTeam(c *gin.Context) {
	var teams []pojo.Team
	offset := c.PostForm("offset")
	offsetInt, _ := strconv.Atoi(offset)
	userId, _ := c.Get("userId")

	if err := config.MysqlDataBase.
		Joins("JOIN join_requests ON join_requests.team_id = teams.id").
		Where("join_requests.user_id = ? AND join_requests.status = ?", userId, 1).
		Offset(offsetInt).
		Limit(10).
		Find(&teams).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队信息时发生错误！"))
		return
	}
	res := teamBodyConvertor(teams)
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
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
	//userId, _ := c.Get("userId")
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
	// userIdUint := userId.(int)
	// if team.LeaderId != uint(userIdUint) {
	// 	c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改该团队的信息"))
	// 	return
	// }
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
		Preload("User").
		Find(&requests).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队申请时发生错误！"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(requests))
}

// 邀请码请求体结构
type JoinByInviteCodeRequest struct {
	InviteCode string `json:"invite_code"`
	UserId     uint   `json:"user_id"`
}

// 通过邀请码加入团队
func JoinTeamByInviteCode(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	var reqBody JoinByInviteCodeRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 查找对应邀请码的团队
	var team pojo.Team
	if err := config.MysqlDataBase.Where("invite_code = ?", reqBody.InviteCode).First(&team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "邀请码无效或团队不存在"))
		return
	}

	// 检查用户是否已经是团队成员
	var existingRequest pojo.JoinRequest
	result := config.MysqlDataBase.Where("user_id = ? AND team_id = ? AND status = ?", userIdInt, team.ID, StatusApproved).
		First(&existingRequest)

	if result.Error == nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "您已经是该团队成员"))
		return
	}

	// 检查用户是否是团队创建者
	if team.LeaderId == uint(userIdInt) {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "您是该团队的创建者，无需加入"))
		return
	}

	// 检查是否有待处理的请求
	var pendingRequest pojo.JoinRequest
	result = config.MysqlDataBase.Where("user_id = ? AND team_id = ? AND status = ?", userIdInt, team.ID, StatusPending).
		First(&pendingRequest)

	if result.Error == nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "您已经提交过加入请求，请等待团队管理员审核"))
		return
	}

	// 创建加入请求
	joinRequest := pojo.JoinRequest{
		UserId: uint(userIdInt),
		TeamId: team.ID,
		Status: StatusPending,
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&joinRequest).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "提交加入请求时发生错误："+err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "提交加入请求时发生错误："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("加入请求已提交", "请等待团队管理员审核您的请求"))
}

// 获取团队成员列表
func GetTeamMembers(c *gin.Context) {
	teamId := c.Query("team_id")
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	// 验证团队存在
	var team pojo.Team
	if err := config.MysqlDataBase.First(&team, teamId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "团队不存在"))
		return
	}

	// 验证用户是否有权限查看（团队成员或创建者）
	if team.LeaderId != uint(userIdInt) {
		var isMember bool
		err := config.MysqlDataBase.Model(&pojo.JoinRequest{}).
			Where("user_id = ? AND team_id = ? AND status = ?", userIdInt, teamId, StatusApproved).
			Select("COUNT(*) > 0").
			Scan(&isMember).Error

		if err != nil || !isMember {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "您没有权限查看该团队成员"))
			return
		}
	}

	// 获取团队创建者信息
	var leader pojo.User
	if err := config.MysqlDataBase.First(&leader, team.LeaderId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队创建者信息失败"))
		return
	}

	// 获取已批准的团队成员
	var members []pojo.User
	if err := config.MysqlDataBase.Model(&pojo.User{}).
		Joins("JOIN join_requests ON join_requests.user_id = users.id").
		Where("join_requests.team_id = ? AND join_requests.status = ?", teamId, StatusApproved).
		Find(&members).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队成员失败"))
		return
	}

	// 构建响应
	type MemberInfo struct {
		User     pojo.User `json:"user"`
		IsLeader bool      `json:"is_leader"`
	}

	var result []MemberInfo

	// 添加团队创建者
	result = append(result, MemberInfo{
		User:     leader,
		IsLeader: true,
	})

	// 添加其他成员
	for _, member := range members {
		result = append(result, MemberInfo{
			User:     member,
			IsLeader: false,
		})
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// 重新生成团队邀请码
func RegenerateInviteCode(c *gin.Context) {
	teamId := c.Query("team_id")
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	// 验证团队存在
	var team pojo.Team
	if err := config.MysqlDataBase.First(&team, teamId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "团队不存在"))
		return
	}

	// 验证用户是否是团队创建者
	if team.LeaderId != uint(userIdInt) {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "只有团队创建者可以重新生成邀请码"))
		return
	}

	// 生成新的邀请码
	newInviteCode := util.GenerateRandomString(8)

	// 更新团队邀请码
	tx := config.MysqlDataBase.Begin()
	team.InviteCode = newInviteCode

	if err := tx.Save(&team).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新邀请码失败："+err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新邀请码失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage(newInviteCode, "邀请码已重新生成"))
}

// 获取团队详情
func GetTeamDetail(c *gin.Context) {
	teamId := c.Query("team_id")
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	// 验证团队存在
	var team pojo.Team
	if err := config.MysqlDataBase.First(&team, teamId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "团队不存在"))
		return
	}

	// 获取团队成员数量
	var memberCount int64
	config.MysqlDataBase.Model(&pojo.JoinRequest{}).
		Where("team_id = ? AND status = ?", team.ID, StatusApproved).
		Count(&memberCount)

	// 检查用户是否是团队创建者
	isLeader := team.LeaderId == uint(userIdInt)

	// 检查用户是否是团队成员
	var isMember bool
	if !isLeader {
		err := config.MysqlDataBase.Model(&pojo.JoinRequest{}).
			Where("user_id = ? AND team_id = ? AND status = ?", userIdInt, team.ID, StatusApproved).
			Select("COUNT(*) > 0").
			Scan(&isMember).Error

		if err != nil {
			isMember = false
		}
	}

	// 构建响应
	type TeamDetailResponse struct {
		Team        pojo.Team `json:"team"`
		MemberCount int64     `json:"member_count"`
		IsLeader    bool      `json:"is_leader"`
		IsMember    bool      `json:"is_member"`
		InviteCode  string    `json:"invite_code,omitempty"` // 只有团队创建者可以看到
	}

	response := TeamDetailResponse{
		Team:        team,
		MemberCount: memberCount + 1, // +1 for the leader
		IsLeader:    isLeader,
		IsMember:    isMember || isLeader,
	}

	// 只有团队创建者可以看到邀请码
	if isLeader {
		response.InviteCode = team.InviteCode
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// 通过邀请码获取团队信息
func GetTeamByInviteCode(c *gin.Context) {
	inviteCode := c.Query("invite_code")

	// 查找对应邀请码的团队
	var team pojo.Team
	if err := config.MysqlDataBase.Where("invite_code = ?", inviteCode).First(&team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "邀请码无效或团队不存在"))
		return
	}

	// 获取团队成员数量
	var memberCount int64
	config.MysqlDataBase.Model(&pojo.JoinRequest{}).
		Where("team_id = ? AND status = ?", team.ID, StatusApproved).
		Count(&memberCount)

	// 构建响应
	type TeamInfoResponse struct {
		Team        pojo.Team `json:"team"`
		MemberCount int64     `json:"memberCount"`
	}

	response := TeamInfoResponse{
		Team:        team,
		MemberCount: memberCount + 1, // +1 表示包含团队创建者
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}
