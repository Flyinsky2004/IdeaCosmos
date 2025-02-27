package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddCreatorComment 添加创作者对版本的评论
func AddCreatorComment(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	var reqBody struct {
		VersionId uint   `json:"version_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 检查版本是否存在
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.First(&version, reqBody.VersionId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "版本不存在"))
		return
	}

	// 检查用户是否有权限评论（项目成员或团队成员）
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.First(&chapter, version.ChapterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "项目不存在"))
		return
	}

	// 检查用户是否是项目成员
	isMember := false

	// 如果是项目创建者
	if project.Team.LeaderId == uint(userIdInt) {
		isMember = true
	}

	// 如果项目有团队，检查用户是否是团队成员
	if !isMember && project.TeamID > 0 {
		var count int64
		config.MysqlDataBase.Model(&pojo.JoinRequest{}).
			Where("team_id = ? AND user_id = ? AND status = ?", project.TeamID, userIdInt, StatusApproved).
			Count(&count)

		if count > 0 {
			isMember = true
		}

		// 检查是否是团队领导
		var team pojo.Team
		if err := config.MysqlDataBase.First(&team, project.TeamID).Error; err == nil {
			if team.LeaderId == uint(userIdInt) {
				isMember = true
			}
		}
	}

	if !isMember {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "您没有权限对此版本进行评论"))
		return
	}

	// 创建评论
	creatorComment := pojo.CreatorComment{
		UserId:    uint(userIdInt),
		VersionId: reqBody.VersionId,
		Content:   reqBody.Content,
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&creatorComment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建评论失败："+err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存评论失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("评论成功", creatorComment))
}

// GetCreatorComments 获取版本的创作者评论列表
func GetCreatorComments(c *gin.Context) {
	versionId := c.Query("version_id")
	versionIdInt, err := strconv.Atoi(versionId)

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的版本ID"))
		return
	}

	var comments []pojo.CreatorComment

	if err := config.MysqlDataBase.Where("version_id = ?", versionIdInt).
		Preload("User").
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取评论失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(comments))
}

// DeleteCreatorComment 删除创作者评论
func DeleteCreatorComment(c *gin.Context) {
	userId, _ := c.Get("userId")
	userIdInt := userId.(int)

	var reqBody struct {
		CommentId uint `json:"comment_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	var comment pojo.CreatorComment
	if err := config.MysqlDataBase.First(&comment, reqBody.CommentId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "评论不存在"))
		return
	}

	// 只有评论作者才能删除评论
	if comment.UserId != uint(userIdInt) {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "您没有权限删除此评论"))
		return
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Delete(&comment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除评论失败："+err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除评论失败："+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("评论已删除"))
}
