package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetIndexCoverList(c *gin.Context) {
	pageIndex := c.Query("pageIndex")
	pageIndexInt, _ := strconv.Atoi(pageIndex)
	var results []pojo.Project
	err := config.MysqlDataBase.Preload("Team").Limit(10).Order("created_at desc").Offset(pageIndexInt).Find(&results).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "查询数据库时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(results))
}

// 获取项目详情
func GetProjectDetail(c *gin.Context) {
	projectId := c.Query("project_id")
	userId := c.Query("user_id") // 从查询参数获取可选的userId

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

	// 增加浏览量并记录观看记录
	tx := config.MysqlDataBase.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新项目观看数
	if err := tx.Model(&project).Update("watches", project.Watches+1).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新观看数失败"))
		return
	}

	// 如果提供了userId,添加观看记录
	if userId != "" {
		uid, err := strconv.Atoi(userId)
		if err == nil {
			watch := pojo.Watch{
				UserId:    uint(uid),
				ProjectId: uint(id),
			}

			if err := tx.Create(&watch).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "记录观看历史失败"))
				return
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "数据库事务提交失败"))
		return
	}

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

// GetChapterDetail 获取篇章详情
func GetChapterDetail(c *gin.Context) {
	chapterId := c.Query("id")

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.
		Preload("CurrentVersion").
		Preload("CurrentVersion.User").
		Where("id = ?", chapterId).
		First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 获取项目信息
	var project pojo.Project
	if err := config.MysqlDataBase.
		Where("id = ?", chapter.ProjectID).
		First(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	// 构建返回数据
	response := struct {
		Chapter pojo.Chapter `json:"chapter"`
		Project pojo.Project `json:"project"`
	}{
		Chapter: chapter,
		Project: project,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// GetHotProjects 获取观看数前10的项目
func GetHotProjects(c *gin.Context) {
	var projects []pojo.Project

	// 直接按观看数降序排序获取前10个项目，预加载Team信息
	err := config.MysqlDataBase.
		Preload("Team").
		Order("watches desc").
		Limit(10).
		Find(&projects).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[[]pojo.Project](500, "获取热门项目失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(projects))
}

// 添加收藏
func AddFavorite(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	// 检查是否已经收藏
	var count int64
	config.MysqlDataBase.Model(&pojo.Favourite{}).
		Where("user_id = ? AND project_id = ?", userId, id).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "已经收藏过该项目"))
		return
	}

	// 创建收藏记录
	favourite := pojo.Favourite{
		UserId:    uint(userId.(int)),
		ProjectId: uint(id),
	}

	err = config.MysqlDataBase.Create(&favourite).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "收藏失败"))
		return
	}

	// 更新项目收藏数
	config.MysqlDataBase.Model(&pojo.Project{}).
		Where("id = ?", id).
		UpdateColumn("favorites", gorm.Expr("favorites + ?", 1))

	c.JSON(http.StatusOK, dto.SuccessResponse("收藏成功"))
}

// 取消收藏
func RemoveFavorite(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	// 删除收藏记录
	result := config.MysqlDataBase.Where("user_id = ? AND project_id = ?", userId, id).
		Delete(&pojo.Favourite{})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "未收藏该项目"))
		return
	}

	// 更新项目收藏数
	config.MysqlDataBase.Model(&pojo.Project{}).
		Where("id = ?", id).
		UpdateColumn("favorites", gorm.Expr("favorites - ?", 1))

	c.JSON(http.StatusOK, dto.SuccessResponse("取消收藏成功"))
}

// 检查是否已收藏
func CheckFavorite(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.SuccessResponse(false))
		return
	}

	projectId := c.Query("project_id")
	id, err := strconv.Atoi(projectId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的项目ID"))
		return
	}

	var count int64
	config.MysqlDataBase.Model(&pojo.Favourite{}).
		Where("user_id = ? AND project_id = ?", userId, id).
		Count(&count)

	c.JSON(http.StatusOK, dto.SuccessResponse(count > 0))
}

// GetWatchesAndLikesAnalysis 获取项目观看和收藏数据分析
func GetWatchesAndLikesAnalysis(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	// 获取用户所属的团队
	var teams []pojo.Team
	if err := config.MysqlDataBase.
		Where("leader_id = ?", userId).
		Or("id IN (SELECT team_id FROM join_requests WHERE user_id = ? AND status = 1)", userId).
		Find(&teams).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队信息失败"))
		return
	}

	teamIds := make([]uint, len(teams))
	for i, team := range teams {
		teamIds[i] = team.ID
	}

	// 获取团队的所有项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.
		Where("team_id IN ?", teamIds).
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	type DailyStats struct {
		Date          string `json:"date"`
		ProjectID     uint   `json:"project_id"`
		ProjectName   string `json:"project_name"`
		WatchCount    int    `json:"watch_count"`
		FavoriteCount int    `json:"favorite_count"`
	}

	var result []DailyStats

	for _, project := range projects {
		// 获取观看数据(前10天)
		var watchStats []struct {
			Date  string
			Count int
		}
		watchQuery := `
			SELECT DATE(created_at) as date, COUNT(*) as count 
			FROM watches 
			WHERE project_id = ? AND deleted_at IS NULL
			GROUP BY DATE(created_at)
			ORDER BY date DESC
			LIMIT 10
		`
		if err := config.MysqlDataBase.Raw(watchQuery, project.ID).Scan(&watchStats).Error; err != nil {
			continue
		}

		// 获取收藏数据(前10天)
		var favoriteStats []struct {
			Date  string
			Count int
		}
		favoriteQuery := `
			SELECT DATE(created_at) as date, 
				   COUNT(CASE WHEN deleted_at IS NULL THEN 1 END) - 
				   COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END) as count
			FROM favourites 
			WHERE project_id = ?
			GROUP BY DATE(created_at)
			ORDER BY date DESC
			LIMIT 10
		`
		if err := config.MysqlDataBase.Raw(favoriteQuery, project.ID).Scan(&favoriteStats).Error; err != nil {
			continue
		}

		// 合并数据
		dateMap := make(map[string]*DailyStats)

		// 处理观看数据
		for _, ws := range watchStats {
			dateMap[ws.Date] = &DailyStats{
				Date:        ws.Date,
				ProjectID:   project.ID,
				ProjectName: project.ProjectName,
				WatchCount:  ws.Count,
			}
		}

		// 处理收藏数据
		for _, fs := range favoriteStats {
			if stats, exists := dateMap[fs.Date]; exists {
				stats.FavoriteCount = fs.Count
			} else {
				dateMap[fs.Date] = &DailyStats{
					Date:          fs.Date,
					ProjectID:     project.ID,
					ProjectName:   project.ProjectName,
					FavoriteCount: fs.Count,
				}
			}
		}

		// 转换为切片
		for _, stats := range dateMap {
			result = append(result, *stats)
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetStyleAndTypeAnalysis 获取不同风格和类型的数据分析
func GetStyleAndTypeAnalysis(c *gin.Context) {
	_, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	// 获取所有项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.Find(&projects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	type StyleStats struct {
		Date       string `json:"date"`
		Style      string `json:"style"`
		Type       string `json:"type"`
		WatchCount int    `json:"watch_count"`
		LikeCount  int    `json:"like_count"`
	}

	var result []StyleStats

	// 获取最近10天的数据
	for _, project := range projects {
		// 解析项目风格
		var styles []string
		if err := json.Unmarshal([]byte(project.Style), &styles); err != nil {
			continue
		}

		// 对每个风格进行统计
		for _, style := range styles {
			// 统计观看数据
			var watchStats []struct {
				Date  string
				Count int
			}
			watchQuery := `
				SELECT DATE(created_at) as date, COUNT(*) as count 
				FROM watches 
				WHERE project_id = ? AND deleted_at IS NULL
				AND created_at >= DATE_SUB(CURDATE(), INTERVAL 10 DAY)
				GROUP BY DATE(created_at)
				ORDER BY date ASC
			`
			if err := config.MysqlDataBase.Raw(watchQuery, project.ID).Scan(&watchStats).Error; err != nil {
				continue
			}

			// 统计收藏数据
			var likeStats []struct {
				Date  string
				Count int
			}
			likeQuery := `
				SELECT DATE(created_at) as date, 
					   COUNT(CASE WHEN deleted_at IS NULL THEN 1 END) - 
					   COUNT(CASE WHEN deleted_at IS NOT NULL THEN 1 END) as count
				FROM favourites 
				WHERE project_id = ?
				AND created_at >= DATE_SUB(CURDATE(), INTERVAL 10 DAY)
				GROUP BY DATE(created_at)
				ORDER BY date ASC
			`
			if err := config.MysqlDataBase.Raw(likeQuery, project.ID).Scan(&likeStats).Error; err != nil {
				continue
			}

			// 获取最近10天的日期范围
			var dates []string
			dateQuery := `
				SELECT DATE(date) as date
				FROM (
					SELECT CURDATE() - INTERVAL (a.a + (10 * b.a) + (100 * c.a)) DAY as date
					FROM (SELECT 0 as a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) as a
					CROSS JOIN (SELECT 0 as a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) as b
					CROSS JOIN (SELECT 0 as a UNION ALL SELECT 1 UNION ALL SELECT 2 UNION ALL SELECT 3 UNION ALL SELECT 4 UNION ALL SELECT 5 UNION ALL SELECT 6 UNION ALL SELECT 7 UNION ALL SELECT 8 UNION ALL SELECT 9) as c
				) dates
				WHERE date >= DATE_SUB(CURDATE(), INTERVAL 10 DAY)
				ORDER BY date ASC
			`
			if err := config.MysqlDataBase.Raw(dateQuery).Scan(&dates).Error; err != nil {
				continue
			}

			// 为每一天创建记录
			for _, date := range dates {
				stats := &StyleStats{
					Date:       date,
					Style:      style,
					Type:       project.Types,
					WatchCount: 0,
					LikeCount:  0,
				}

				// 填充观看数据
				for _, ws := range watchStats {
					if ws.Date == date {
						stats.WatchCount = ws.Count
						break
					}
				}

				// 填充收藏数据
				for _, ls := range likeStats {
					if ls.Date == date {
						stats.LikeCount = ls.Count
						break
					}
				}

				result = append(result, *stats)
			}
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// AddVersionFeeling 添加用户对版本的情绪评价
func AddVersionFeeling(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	var feeling pojo.Feeling
	if err := c.ShouldBindJSON(&feeling); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	feeling.UserId = uint(userId.(int))

	// 检查是否已经评价过
	var count int64
	config.MysqlDataBase.Model(&pojo.Feeling{}).
		Where("user_id = ? AND version_id = ?", feeling.UserId, feeling.VersionId).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "已经评价过该版本"))
		return
	}

	// 创建新的情绪评价
	if err := config.MysqlDataBase.Create(&feeling).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "添加情绪评价失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("评价成功"))
}

// GetVersionFeeling 获取用户对版本的情绪评价
func GetVersionFeeling(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "未登录"))
		return
	}

	versionId := c.Query("version_id")
	vid, err := strconv.Atoi(versionId)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的版本ID"))
		return
	}

	var feeling pojo.Feeling
	err = config.MysqlDataBase.Where("user_id = ? AND version_id = ?", userId, vid).First(&feeling).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.SuccessResponse("获取失败"))
			return
		}
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取情绪评价失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(feeling))
}

// GetEmotionAnalysis 获取情绪分析数据
func GetEmotionAnalysis(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	// 获取用户所属的团队
	var teams []pojo.Team
	if err := config.MysqlDataBase.
		Where("leader_id = ?", userId).
		Or("id IN (SELECT team_id FROM join_requests WHERE user_id = ? AND status = 1)", userId).
		Find(&teams).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取团队信息失败"))
		return
	}

	teamIds := make([]uint, len(teams))
	for i, team := range teams {
		teamIds[i] = team.ID
	}

	// 获取团队的所有项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.
		Where("team_id IN ?", teamIds).
		Find(&projects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	type EmotionStats struct {
		Date        string `json:"date"`
		ProjectID   uint   `json:"project_id"`
		ProjectName string `json:"project_name"`
		Emotion     string `json:"emotion"`
		Count       int    `json:"count"`
	}

	var result []EmotionStats

	for _, project := range projects {
		// 获取项目所有章节
		var chapters []pojo.Chapter
		if err := config.MysqlDataBase.
			Where("project_id = ?", project.ID).
			Preload("CurrentVersion").
			Find(&chapters).Error; err != nil {
			continue
		}

		// 获取所有版本ID
		versionIds := make([]uint, len(chapters))
		for i, chapter := range chapters {
			if chapter.CurrentVersion.ID != 0 {
				versionIds[i] = chapter.CurrentVersion.ID
			}
		}

		// 获取情绪统计数据
		var stats []struct {
			Date    string
			Emotion string
			Count   int
		}

		query := `
			SELECT 
				DATE(created_at) as date,
				feeling as emotion,
				COUNT(*) as count
			FROM feelings
			WHERE version_id IN ?
				AND created_at >= DATE_SUB(CURDATE(), INTERVAL 10 DAY)
				AND deleted_at IS NULL
			GROUP BY DATE(created_at), feeling
			ORDER BY date DESC
		`

		if err := config.MysqlDataBase.Raw(query, versionIds).Scan(&stats).Error; err != nil {
			continue
		}

		// 转换为响应格式
		for _, stat := range stats {
			result = append(result, EmotionStats{
				Date:        stat.Date,
				ProjectID:   project.ID,
				ProjectName: project.ProjectName,
				Emotion:     stat.Emotion,
				Count:       stat.Count,
			})
		}
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}
