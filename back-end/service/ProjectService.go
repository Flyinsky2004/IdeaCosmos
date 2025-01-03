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
 * @date 2024/12/29 20:43
 */
//获取某个用户创建或加入的项目
func GetProjectsByUserId(userId uint) ([]pojo.Project, error) {
	var projects []pojo.Project

	// 子查询：用户创建的团队
	createdTeamsSubQuery := config.MysqlDataBase.Model(&pojo.Team{}).Select("id").Where("leader_id = ?", userId)

	// 子查询：用户加入的团队
	joinedTeamsSubQuery := config.MysqlDataBase.Model(&pojo.JoinRequest{}).
		Select("team_id").
		Where("user_id = ? AND status = ?", userId, 1)

	// 查询项目：团队 ID 在上述两个子查询结果中
	err := config.MysqlDataBase.Preload("Team"). // 预加载 Team 信息
							Where("team_id IN (?) OR team_id IN (?)", createdTeamsSubQuery, joinedTeamsSubQuery).
							Find(&projects).Error

	return projects, err
}

// 获取项目列表
func GetProjectList(c *gin.Context) {
	userId, _ := c.Get("userId")
	prsm, err := GetProjectsByUserId(uint(userId.(int)))
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询数据库时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage[[]pojo.Project]("查询成功", prsm))
}

// CreateProject 创建新项目
func CreateProject(c *gin.Context) {
	var project pojo.Project
	userId, _ := c.Get("userId")
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误,错误:"+err.Error()))
		return
	}
	var TeamRequest pojo.JoinRequest
	var Team pojo.Team
	if err := config.MysqlDataBase.Where("id = ?", project.TeamID).First(&Team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "团队不存在"))
		return
	}
	if Team.LeaderId != uint(userId.(int)) {
		if err := config.MysqlDataBase.Where("team_id = ? AND user_id = ? AND status = 1", project.TeamID, userId).First(&TeamRequest).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您尚未加入其项目工作团队"))
			return
		}
	}
	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建项目失败 code: 1"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建项目失败 code: 2"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage[string]("项目创建成功！", ""))
}

// UpdateProject 更新项目信息
func UpdateProject(c *gin.Context) {
	var project pojo.Project
	userId, _ := c.Get("userId")

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}
	var projectSource pojo.Project
	if err := config.MysqlDataBase.Where("ID = ?", project.ID).First(&projectSource).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "项目不存在"))
		return
	}
	var TeamRequest pojo.JoinRequest
	var Team pojo.Team
	if err := config.MysqlDataBase.Where("id = ?", projectSource.TeamID).First(&Team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "团队不存在"))
		return
	}
	if Team.LeaderId != uint(userId.(int)) {
		if err := config.MysqlDataBase.Where("team_id = ? AND user_id = ? AND status = 1", projectSource.TeamID, userId).First(&TeamRequest).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您尚未加入其项目工作团队"))
			return
		}
	}

	projectSource.ProjectName = project.ProjectName
	projectSource.CoverImage = project.CoverImage
	projectSource.CustomPrompt = project.CustomPrompt
	projectSource.HighPoint = project.HighPoint
	projectSource.Start = project.Start
	projectSource.Resolved = project.Resolved
	projectSource.SocialStory = project.SocialStory
	projectSource.Style = project.Style
	projectSource.MarketPeople = project.MarketPeople
	projectSource.Types = project.Types

	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&projectSource).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存项目信息时发生错误。code : 1"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存项目信息时发生错误。code : 2"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[string]("更新成功"))
}

// CreateCharacter 创建角色
func CreateCharacter(c *gin.Context) {
	userId, _ := c.Get("userId")
	var character pojo.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}
	var projectSource pojo.Project
	if err := config.MysqlDataBase.Where("ID = ?", character.ProjectID).First(&projectSource).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "项目不存在"))
		return
	}
	var TeamRequest pojo.JoinRequest
	var Team pojo.Team
	if err := config.MysqlDataBase.Where("id = ?", projectSource.TeamID).First(&Team).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "团队不存在"))
		return
	}
	if Team.LeaderId != uint(userId.(int)) {
		if err := config.MysqlDataBase.Where("team_id = ? AND user_id = ? AND status = 1", projectSource.TeamID, userId).First(&TeamRequest).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您尚未加入其项目工作团队"))
			return
		}
	}
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&character).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存角色信息时发生错误。code : 1"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存角色信息时发生错误。code : 2"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage[string]("角色添加成功！", ""))
}

// 检查用户是否有权限操作项目
func checkProjectPermission(userId uint, projectID uint) (bool, error) {
	var project pojo.Project
	if err := config.MysqlDataBase.Where("ID = ?", projectID).First(&project).Error; err != nil {
		return false, err
	}

	var team pojo.Team
	if err := config.MysqlDataBase.Where("id = ?", project.TeamID).First(&team).Error; err != nil {
		return false, err
	}

	if team.LeaderId == userId {
		return true, nil
	}

	var teamRequest pojo.JoinRequest
	if err := config.MysqlDataBase.Where("team_id = ? AND user_id = ? AND status = 1",
		project.TeamID, userId).First(&teamRequest).Error; err != nil {
		return false, err
	}

	return true, nil
}

// UpdateCharacter 更新角色信息
func UpdateCharacter(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
		return
	}

	var character pojo.Character
	// 绑定更新数据
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}
	// 获取原有角色信息
	if err := config.MysqlDataBase.First(&character, character.ID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "角色不存在"))
		return
	}

	// 检查权限
	hasPermission, err := checkProjectPermission(uint(userId.(int)), character.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
		return
	}
	if !hasPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改该角色"))
		return
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Model(&pojo.Character{}).Where("id = ?", character.ID).Updates(character).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[string]("更新成功"))
}

// CreateCharacterRelationship 创建角色关系
//func CreateCharacterRelationship(c *gin.Context) {
//	userId, exists := c.Get("userId")
//	if !exists {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
//		return
//	}
//
//	var relationship pojo.CharacterRelationShip
//	if err := c.ShouldBindJSON(&relationship); err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
//		return
//	}
//
//	// 验证两个角色是否存在并获取项目ID
//	var firstChar, secondChar pojo.Character
//	if err := config.MysqlDataBase.First(&firstChar, relationship.FirstCharacterID).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "第一个角色不存在"))
//		return
//	}
//	if err := config.MysqlDataBase.First(&secondChar, relationship.SecondCharacterID).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "第二个角色不存在"))
//		return
//	}
//
//	// 验证两个角色是否属于同一个项目
//	if firstChar.ProjectID != secondChar.ProjectID {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "两个角色必须属于同一个项目"))
//		return
//	}
//
//	// 检查权限
//	hasPermission, err := checkProjectPermission(uint(userId.(int)), firstChar.ProjectID)
//	if err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
//		return
//	}
//	if !hasPermission {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限创建角色关系"))
//		return
//	}
//
//	tx := config.MysqlDataBase.Begin()
//	if err := tx.Create(&relationship).Error; err != nil {
//		tx.Rollback()
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建角色关系失败"))
//		return
//	}
//	if err := tx.Commit().Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建角色关系失败"))
//		return
//	}
//
//	c.JSON(http.StatusOK, dto.SuccessResponse[pojo.CharacterRelationShip](relationship))
//}
//
//// UpdateCharacterRelationship 更新角色关系
//func UpdateCharacterRelationship(c *gin.Context) {
//	userId, exists := c.Get("userId")
//	if !exists {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
//		return
//	}
//
//	var relationship pojo.CharacterRelationShip
//	relationshipId := c.Param("id")
//
//	// 获取原有关系信息
//	if err := config.MysqlDataBase.First(&relationship, relationshipId).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "角色关系不存在"))
//		return
//	}
//
//	// 获取角色信息以验证权限
//	var character pojo.Character
//	if err := config.MysqlDataBase.First(&character, relationship.FirstCharacterID).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色信息失败"))
//		return
//	}
//
//	// 检查权限
//	hasPermission, err := checkProjectPermission(uint(userId.(int)), character.ProjectID)
//	if err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
//		return
//	}
//	if !hasPermission {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改角色关系"))
//		return
//	}
//
//	if err := c.ShouldBindJSON(&relationship); err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
//		return
//	}
//
//	tx := config.MysqlDataBase.Begin()
//	if err := tx.Model(&pojo.CharacterRelationShip{}).
//		Where("id = ?", relationshipId).
//		Updates(relationship).Error; err != nil {
//		tx.Rollback()
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色关系失败"))
//		return
//	}
//	if err := tx.Commit().Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色关系失败"))
//		return
//	}
//
//	c.JSON(http.StatusOK, dto.SuccessResponse[string]("更新成功"))
//}
//
//// GetProjectCharacters 获取项目下的所有角色
//func GetProjectCharacters(c *gin.Context) {
//	userId, exists := c.Get("userId")
//	if !exists {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
//		return
//	}
//
//	projectId := c.Param("id")
//
//	// 检查权限
//	hasPermission, err := checkProjectPermission(uint(userId.(int)), uint(projectId))
//	if err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
//		return
//	}
//	if !hasPermission {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限查看该项目的角色"))
//		return
//	}
//
//	var characters []pojo.Character
//	if err := config.MysqlDataBase.Where("project_id = ?", projectId).
//		Find(&characters).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色列表失败"))
//		return
//	}
//
//	c.JSON(http.StatusOK, dto.SuccessResponse[[]pojo.Character](characters))
//}
//
//// GetCharacterRelationships 获取角色的所有关系
//func GetCharacterRelationships(c *gin.Context) {
//	userId, exists := c.Get("userId")
//	if !exists {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
//		return
//	}
//
//	characterId := c.Param("id")
//
//	// 获取角色信息以验证权限
//	var character pojo.Character
//	if err := config.MysqlDataBase.First(&character, characterId).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "角色不存在"))
//		return
//	}
//
//	// 检查权限
//	hasPermission, err := checkProjectPermission(uint(userId.(int)), character.ProjectID)
//	if err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
//		return
//	}
//	if !hasPermission {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限查看该角色的关系"))
//		return
//	}
//
//	var relationships []pojo.CharacterRelationShip
//	if err := config.MysqlDataBase.
//		Preload("FirstCharacter").
//		Preload("SecondCharacter").
//		Where("first_character_id = ? OR second_character_id = ?", characterId, characterId).
//		Find(&relationships).Error; err != nil {
//		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色关系失败"))
//		return
//	}
//
//	c.JSON(http.StatusOK, dto.SuccessResponse[[]pojo.CharacterRelationShip](relationships))
//}
