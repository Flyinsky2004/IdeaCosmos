package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIndexCoverList(c *gin.Context) {
	pageIndex := c.Query("pageIndex")
	pageIndexInt, _ := strconv.Atoi(pageIndex)
	var results []pojo.Project
	err := config.MysqlDataBase.Preload("Team").Limit(10).Offset(pageIndexInt).Find(&results).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询数据库时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(results))
}

// 获取项目详情
func GetProjectDetail(c *gin.Context) {
	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	var project pojo.Project
	// 预加载 Team 信息并获取项目
	err = config.MysqlDataBase.Preload("Team").First(&project, id).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "项目不存在"))
		return
	}

	// 增加浏览量
	config.MysqlDataBase.Model(&project).Update("watches", project.Watches+1)

	c.JSON(http.StatusOK, dto.SuccessResponse(project))
}

// 获取项目角色列表
func GetProjectCharacters(c *gin.Context) {
	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	var characters []pojo.Character
	err = config.MysqlDataBase.Where("project_id = ?", id).Find(&characters).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询角色信息失败"))
		return
	}

	// 获取角色关系
	var relationships []pojo.CharacterRelationShip
	err = config.MysqlDataBase.Where("first_character_id IN (?) OR second_character_id IN (?)",
		getCharacterIds(characters), getCharacterIds(characters)).
		Preload("FirstCharacter").
		Preload("SecondCharacter").
		Find(&relationships).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询角色关系失败"))
		return
	}

	response := struct {
		Characters    []pojo.Character             `json:"characters"`
		Relationships []pojo.CharacterRelationShip `json:"relationships"`
	}{
		Characters:    characters,
		Relationships: relationships,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// 获取项目章节列表
func GetProjectChapters(c *gin.Context) {
	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	var chapters []pojo.Chapter
	err = config.MysqlDataBase.Where("project_id = ?", id).
		Preload("CurrentVersion").
		Preload("CurrentVersion.User").
		Find(&chapters).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询章节信息失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(chapters))
}

// 辅助函数：获取角色ID列表
func getCharacterIds(characters []pojo.Character) []uint {
	ids := make([]uint, len(characters))
	for i, char := range characters {
		ids[i] = char.ID
	}
	return ids
}

// 获取项目评论
func GetProjectComments(c *gin.Context) {
	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	var comments []pojo.ProjectComment
	err = config.MysqlDataBase.Where("project_id = ?", id).
		Preload("User").
		Order("created_at desc").
		Find(&comments).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取评论失败"+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(comments))
}

type ProjectCommentRequest struct {
	Content   string `json:"content"`
	ProjectId int    `json:"project_id"`
}

// 添加项目评论
func AddProjectComment(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	var commentRequest ProjectCommentRequest
	if err := c.ShouldBindJSON(&commentRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}
	projectIdUint := uint(commentRequest.ProjectId)
	comment := pojo.ProjectComment{
		Content:   commentRequest.Content,
		ProjectId: projectIdUint,
		UserId:    uint(userId.(int)),
	}

	comment.UserId = uint(userId.(int))
	err := config.MysqlDataBase.Create(&comment).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "添加评论失败"))
		return
	}

	// 返回新创建的评论（包含用户信息）
	config.MysqlDataBase.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusOK, dto.SuccessResponse(comment))
}
