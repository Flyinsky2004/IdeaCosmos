package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetIndexCoverList 获取首页封面列表
func GetIndexCoverList(c *gin.Context) {
	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // 默认获取5条记录
	}

	sortBy := c.Query("sort_by")
	switch sortBy {
	case "new":
		// 获取最新项目
		var projects []pojo.Project
		if err := config.MysqlDataBase.Where("status != ?", "banned").Preload("Team").Order("created_at DESC").Limit(limit).Find(&projects).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse(projects))
		return
	case "hot":
		// 获取热门项目（根据浏览量）
		var projects []pojo.Project
		query := `
		SELECT p.* FROM projects p
		LEFT JOIN (
			SELECT project_id, COUNT(*) as watch_count
			FROM watches
			WHERE watches.deleted_at IS NULL 
			GROUP BY project_id
		) w ON p.id = w.project_id
		WHERE p.deleted_at IS NULL AND p.status != 'banned'
		ORDER BY w.watch_count DESC, p.created_at DESC
		LIMIT ?
		`
		if err := config.MysqlDataBase.Raw(query, limit).Scan(&projects).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
			return
		}

		// 由于Raw查询不支持Preload，需要单独加载Team信息
		if err := config.MysqlDataBase.Preload("Team").Find(&projects, "id IN ?", getProjectIds(projects)).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目团队信息失败"))
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse(projects))
		return
	case "featured":
		// 获取推荐项目
		var projects []pojo.Project
		if err := config.MysqlDataBase.Where("status = ?", "featured").Preload("Team").Order("created_at DESC").Limit(limit).Find(&projects).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse(projects))
		return
	default:
		// 默认随机获取项目
		var projects []pojo.Project
		if err := config.MysqlDataBase.Where("status != ?", "banned").Preload("Team").Order("RAND()").Limit(limit).Find(&projects).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
			return
		}

		c.JSON(http.StatusOK, dto.SuccessResponse(projects))
		return
	}
}

// getProjectIds 辅助函数：从项目列表中获取所有项目ID
func getProjectIds(projects []pojo.Project) []uint {
	ids := make([]uint, len(projects))
	for i, project := range projects {
		ids[i] = project.ID
	}
	return ids
}

// GetProjectDetail 获取项目详情
func GetProjectDetail(c *gin.Context) {
	id := c.Query("id")
	var project pojo.Project
	if err := config.MysqlDataBase.Preload("Team").First(&project, id).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "项目不存在"))
		return
	}

	// 检查项目状态
	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	// 获取项目的浏览量
	var watchCount int64
	if err := config.MysqlDataBase.Model(&pojo.Watch{}).Where("project_id = ?", id).Count(&watchCount).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目浏览量失败"))
		return
	}

	// 获取项目的收藏量
	var favoriteCount int64
	if err := config.MysqlDataBase.Model(&pojo.Favourite{}).Where("project_id = ?", id).Count(&favoriteCount).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目收藏量失败"))
		return
	}

	// 获取项目的章节数
	var chapterCount int64
	if err := config.MysqlDataBase.Model(&pojo.Chapter{}).Where("project_id = ?", id).Count(&chapterCount).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目章节数失败"))
		return
	}

	// 如果用户已登录，检查是否已收藏
	var isFavorite bool
	userId, exists := c.Get("userId")
	if exists {
		var count int64
		config.MysqlDataBase.Model(&pojo.Favourite{}).
			Where("user_id = ? AND project_id = ?", userId, id).
			Count(&count)
		isFavorite = count > 0
	}

	// 收集项目详情
	var projectDetail struct {
		pojo.Project
		WatchCount    int64 `json:"watch_count"`
		FavoriteCount int64 `json:"favorite_count"`
		ChapterCount  int64 `json:"chapter_count"`
		IsFavorite    bool  `json:"is_favorite"`
	}

	projectDetail.Project = project
	projectDetail.WatchCount = watchCount
	projectDetail.FavoriteCount = favoriteCount
	projectDetail.ChapterCount = chapterCount
	projectDetail.IsFavorite = isFavorite

	c.JSON(http.StatusOK, dto.SuccessResponse(projectDetail))
}

// GetProjectCharacters 获取项目角色
func GetProjectCharacters(c *gin.Context) {
	projectId := c.Query("id")

	// 检查项目状态
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "项目不存在"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", projectId).Find(&characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色列表失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(characters))
}

// GetProjectChapters 获取项目章节列表
func GetProjectChapters(c *gin.Context) {
	projectId := c.Query("id")

	// 检查项目状态
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "项目不存在"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	var chapters []pojo.Chapter
	if err := config.MysqlDataBase.Where("project_id = ?", projectId).
		Order("created_at ASC").
		Preload("CurrentVersion", "status = 'approved'"). // 只加载已审核通过的当前版本
		Find(&chapters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节列表失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(chapters))
}

// GetChapterDetail 获取章节详细信息
func GetChapterDetail(c *gin.Context) {
	chapterId := c.Query("id")

	// 获取章节信息，同时预加载当前版本和项目信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.
		Preload("CurrentVersion", "status = 'approved'"). // 只加载已审核通过的当前版本
		First(&chapter, chapterId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "章节不存在"))
		return
	}

	// 获取项目信息
	var project pojo.Project
	if err := config.MysqlDataBase.
		Preload("Team").
		First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	// 检查项目状态
	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	// 检查是否有已审核通过的当前版本
	if chapter.CurrentVersion.ID == 0 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该章节暂无已审核通过的内容"))
		return
	}

	// 构造响应数据
	response := struct {
		Chapter pojo.Chapter `json:"chapter"`
		Project pojo.Project `json:"project"`
	}{
		Chapter: chapter,
		Project: project,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
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

	// 向团队Leader发送通知
	userName := getUserName(uint(userId.(int)))
	notificationTitle := "新的项目评论"
	notificationContent := fmt.Sprintf("用户 %s 对您团队的项目添加了评论", userName)
	sendErr := sendNotificationToTeamLeader(
		pojo.CommentNotification,
		uint(userId.(int)),
		projectIdUint,
		notificationTitle,
		notificationContent,
		comment.ID,
		"project_comment",
	)
	if sendErr != nil {
		// 仅记录错误，不影响主流程
		log.Printf("发送通知失败: %v", sendErr)
	}

	// 返回新创建的评论（包含用户信息）
	config.MysqlDataBase.Preload("User").First(&comment, comment.ID)
	c.JSON(http.StatusOK, dto.SuccessResponse(comment))
}

// GetHotProjects 获取热门项目
func GetHotProjects(c *gin.Context) {
	var hotProjects []struct {
		pojo.Project
		WatchCount int64 `json:"watch_count"`
	}

	// 获取热门项目，根据浏览量排序，排除已下架项目
	query := `
	SELECT p.*, COALESCE(w.watch_count, 0) as watch_count
	FROM projects p
	LEFT JOIN (
		SELECT project_id, COUNT(*) as watch_count
		FROM watches
		WHERE watches.deleted_at IS NULL
		GROUP BY project_id
	) w ON p.id = w.project_id
	WHERE p.deleted_at IS NULL AND p.status != 'banned'
	ORDER BY w.watch_count DESC, p.created_at DESC
	LIMIT 10
	`

	if err := config.MysqlDataBase.Raw(query).Scan(&hotProjects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取热门项目失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(hotProjects))
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

	// 向团队Leader发送通知
	userName := getUserName(uint(userId.(int)))
	notificationTitle := "新的项目收藏"
	notificationContent := fmt.Sprintf("用户 %s 收藏了您团队的项目", userName)

	// 获取项目信息用于通知
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, id).Error; err == nil {
		sendErr := sendNotificationToTeamLeader(
			pojo.LikeNotification,
			uint(userId.(int)),
			uint(id),
			notificationTitle,
			notificationContent,
			favourite.ID,
			"project_favorite",
		)
		if sendErr != nil {
			// 仅记录错误，不影响主流程
			log.Printf("发送通知失败: %v", sendErr)
		}
	}

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

	// 获取团队的所有项目 - 只考虑未被下架的项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.
		Where("team_id IN ? AND status != ?", teamIds, "banned").
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

type StyleStats struct {
	Date       string `json:"date"`
	Style      string `json:"style"`
	Type       string `json:"type"`
	WatchCount int    `json:"watch_count"`
	LikeCount  int    `json:"like_count"`
}

// GetStyleAndTypeAnalysis 获取不同风格和类型的数据分析
func GetStyleAndTypeAnalysis(c *gin.Context) {
	_, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	// 获取所有未被下架的项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.Where("status != ?", "banned").Find(&projects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
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

	// 检查版本状态
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.First(&version, vid).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "版本不存在"))
		return
	}

	if version.Status != "approved" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该章节版本未通过审核"))
		return
	}

	// 检查项目状态
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.First(&chapter, version.ChapterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "章节不存在"))
		return
	}

	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "项目不存在"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
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

	// 获取版本信息并检查状态
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.First(&version, feeling.VersionId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "版本不存在"))
		return
	}

	// 检查版本状态是否已审核通过
	if version.Status != "approved" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该章节版本未通过审核，无法评价"))
		return
	}

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.First(&chapter, version.ChapterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节信息失败"))
		return
	}

	// 检查项目状态
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架，无法评价"))
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

	// 向团队Leader发送通知
	userName := getUserName(uint(userId.(int)))
	notificationTitle := "新的情绪评价"
	notificationContent := fmt.Sprintf("用户 %s 对您团队的项目章节版本添加了情绪评价", userName)
	sendErr := sendNotificationToTeamLeader(
		pojo.ContentUpdateNotification,
		uint(userId.(int)),
		project.ID,
		notificationTitle,
		notificationContent,
		feeling.ID,
		"version_feeling",
	)
	if sendErr != nil {
		// 仅记录错误，不影响主流程
		log.Printf("发送通知失败: %v", sendErr)
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("评价成功"))
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

	// 获取团队的所有项目 - 只考虑未被下架的项目
	var projects []pojo.Project
	if err := config.MysqlDataBase.
		Where("team_id IN ? AND status != ?", teamIds, "banned").
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
			Preload("CurrentVersion", "status = 'approved'"). // 只获取已通过审核的版本
			Find(&chapters).Error; err != nil {
			continue
		}

		// 获取所有通过审核的版本ID
		var versionIds []uint
		for _, chapter := range chapters {
			if chapter.CurrentVersion.ID != 0 {
				versionIds = append(versionIds, chapter.CurrentVersion.ID)
			}
		}

		// 如果没有通过审核的版本，跳过
		if len(versionIds) == 0 {
			continue
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

// GetVersionComments 获取版本评论
func GetVersionComments(c *gin.Context) {
	versionId := c.Query("version_id")

	// 检查版本状态
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.First(&version, versionId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "版本不存在"))
		return
	}

	if version.Status != "approved" {
		log.Println(version.Status)
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该章节版本未通过审核"))
		return
	}

	// 获取章节和项目信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.First(&chapter, version.ChapterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节信息失败"))
		return
	}

	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	// 获取作者评论
	var authorComments []pojo.AuthorComment
	if err := config.MysqlDataBase.
		Where("version_id = ?", versionId).
		Preload("User").
		Find(&authorComments).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取作者评论失败"))
		return
	}

	// 获取读者评论
	var readerComments []pojo.ReaderComment
	if err := config.MysqlDataBase.
		Where("version_id = ?", versionId).
		Preload("User").
		Find(&readerComments).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取读者评论失败"))
		return
	}

	// 响应数据
	var response struct {
		AuthorComments []pojo.AuthorComment `json:"author_comments"`
		ReaderComments []pojo.ReaderComment `json:"reader_comments"`
	}
	response.AuthorComments = authorComments
	response.ReaderComments = readerComments

	c.JSON(http.StatusOK, dto.SuccessResponse(response))
}

// AddVersionComment 添加版本评论
func AddVersionComment(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录"))
		return
	}

	var requestBody struct {
		VersionId uint   `json:"version_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Type      string `json:"type" binding:"required"` // "reader" 或 "author"
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 获取版本信息并检查状态
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.First(&version, requestBody.VersionId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "版本不存在"))
		return
	}

	// 检查版本状态是否已审核通过
	if version.Status != "approved" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该章节版本未通过审核，无法评论"))
		return
	}

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.First(&chapter, version.ChapterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节信息失败"))
		return
	}

	// 检查项目状态
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, chapter.ProjectID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目信息失败"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架，无法评论"))
		return
	}

	// 检查权限
	if requestBody.Type == "author" {
		// 只有项目作者才能添加作者评论
		isTeamMember, err := isUserInProjectTeam(uint(userId.(int)), chapter.ProjectID)
		if err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限失败"))
			return
		}

		if !isTeamMember {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "只有项目团队成员才能添加作者评论"))
			return
		}

		// 创建作者评论
		authorComment := pojo.AuthorComment{
			Content:   requestBody.Content,
			VersionId: requestBody.VersionId,
			UserId:    uint(userId.(int)),
		}

		if err := config.MysqlDataBase.Create(&authorComment).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "添加评论失败"))
			return
		}

		// 向团队Leader发送通知
		userName := getUserName(uint(userId.(int)))
		notificationTitle := "新的作者评论"
		notificationContent := fmt.Sprintf("用户 %s 在您团队的项目章节版本中添加了作者评论", userName)
		sendErr := sendNotificationToTeamLeader(
			pojo.CommentNotification,
			uint(userId.(int)),
			project.ID,
			notificationTitle,
			notificationContent,
			authorComment.ID,
			"author_comment",
		)
		if sendErr != nil {
			// 仅记录错误，不影响主流程
			log.Printf("发送通知失败: %v", sendErr)
		}

	} else if requestBody.Type == "reader" {
		// 创建读者评论
		readerComment := pojo.ReaderComment{
			Content:   requestBody.Content,
			VersionId: requestBody.VersionId,
			UserId:    uint(userId.(int)),
		}

		if err := config.MysqlDataBase.Create(&readerComment).Error; err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "添加评论失败"))
			return
		}

		// 向团队Leader发送通知
		userName := getUserName(uint(userId.(int)))
		notificationTitle := "新的读者评论"
		notificationContent := fmt.Sprintf("用户 %s 在您团队的项目章节版本中添加了读者评论", userName)
		sendErr := sendNotificationToTeamLeader(
			pojo.CommentNotification,
			uint(userId.(int)),
			project.ID,
			notificationTitle,
			notificationContent,
			readerComment.ID,
			"reader_comment",
		)
		if sendErr != nil {
			// 仅记录错误，不影响主流程
			log.Printf("发送通知失败: %v", sendErr)
		}
	} else {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的评论类型"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("评论成功"))
}

// 辅助函数: 检查用户是否是项目团队成员
func isUserInProjectTeam(userId uint, projectId uint) (bool, error) {
	var project pojo.Project
	if err := config.MysqlDataBase.Preload("Team").First(&project, projectId).Error; err != nil {
		return false, err
	}

	if project.Team.LeaderId == userId {
		return true, nil
	}

	var count int64
	err := config.MysqlDataBase.Model(&pojo.JoinRequest{}).
		Where("team_id = ? AND user_id = ? AND status = 1", project.TeamID, userId).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetProjectCharacterRelationships 获取项目下所有角色关系
func GetProjectCharacterRelationships(c *gin.Context) {
	projectId := c.Query("id")

	// 检查项目状态
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, projectId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "项目不存在"))
		return
	}

	if project.Status == "banned" {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "该项目已下架"))
		return
	}

	// 获取项目下的所有角色ID
	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", projectId).Find(&characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色列表失败"))
		return
	}

	// 如果没有角色，直接返回空数组
	if len(characters) == 0 {
		c.JSON(http.StatusOK, dto.SuccessResponse([]pojo.CharacterRelationShip{}))
		return
	}

	// 获取角色ID列表
	characterIds := make([]uint, len(characters))
	for i, char := range characters {
		characterIds[i] = char.ID
	}

	// 获取所有角色关系
	var relationships []pojo.CharacterRelationShip
	if err := config.MysqlDataBase.
		Where("first_character_id IN ? OR second_character_id IN ?", characterIds, characterIds).
		Preload("FirstCharacter").
		Preload("SecondCharacter").
		Find(&relationships).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色关系失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(relationships))
}

// 辅助函数：发送通知给团队Leader
func sendNotificationToTeamLeader(notificationType pojo.NotificationType, senderID uint, projectID uint, title string, content string, relatedID uint, relatedType string) error {
	// 获取项目信息
	var project pojo.Project
	if err := config.MysqlDataBase.First(&project, projectID).Error; err != nil {
		return err
	}

	// 获取团队信息
	var team pojo.Team
	if err := config.MysqlDataBase.First(&team, project.TeamID).Error; err != nil {
		return err
	}

	// 创建通知
	notification := pojo.Notification{
		Type:        notificationType,
		Title:       title,
		Content:     content,
		SenderID:    senderID,
		ReceiverID:  team.LeaderId,
		IsRead:      false,
		ReadTime:    nil,
		RelatedID:   relatedID,
		RelatedType: relatedType,
		ExtraData:   "",
	}

	// 保存通知
	if err := config.MysqlDataBase.Create(&notification).Error; err != nil {
		return err
	}

	return nil
}

// 辅助函数：获取用户名称
func getUserName(userID uint) string {
	var user pojo.User
	if err := config.MysqlDataBase.Select("username").First(&user, userID).Error; err != nil {
		return "未知用户"
	}
	return user.Username
}
