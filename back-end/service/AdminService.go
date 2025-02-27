package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 检查用户是否为管理员
func isAdmin(userId int) bool {
	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userId).First(&user).Error; err != nil {
		return false
	}
	return user.Permission >= 1 // 权限值为1以上表示管理员
}

// AdminAuthMiddleware 管理员权限中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或会话已过期"))
			c.Abort()
			return
		}

		if !isAdmin(userId.(int)) {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "您没有管理员权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}

// GetAdminDashboard 获取管理面板数据
func GetAdminDashboard(c *gin.Context) {
	// 计算各类统计数据
	var userCount int64
	var projectCount int64
	var chapterCount int64
	var commentCount int64

	config.MysqlDataBase.Model(&pojo.User{}).Count(&userCount)
	config.MysqlDataBase.Model(&pojo.Project{}).Count(&projectCount)
	config.MysqlDataBase.Model(&pojo.Chapter{}).Count(&chapterCount)

	// 评论数量需要合并多个表
	var readerCommentCount int64
	var authorCommentCount int64
	var projectCommentCount int64
	config.MysqlDataBase.Model(&pojo.ReaderComment{}).Count(&readerCommentCount)
	config.MysqlDataBase.Model(&pojo.AuthorComment{}).Count(&authorCommentCount)
	config.MysqlDataBase.Model(&pojo.ProjectComment{}).Count(&projectCommentCount)
	commentCount = readerCommentCount + authorCommentCount + projectCommentCount

	// 最近注册的用户
	var recentUsers []pojo.User
	config.MysqlDataBase.Order("created_at desc").Limit(5).Find(&recentUsers)
	for i := range recentUsers {
		recentUsers[i].Password = "" // 移除密码字段
	}

	// 最热门项目
	var hotProjects []pojo.Project
	config.MysqlDataBase.Order("watches desc").Limit(5).Find(&hotProjects)

	// 最近活跃用户（基于评论）
	type ActiveUser struct {
		UserID        uint   `json:"user_id"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		ActivityCount int    `json:"activity_count"`
	}
	var activeUsers []ActiveUser
	config.MysqlDataBase.Raw(`
		SELECT u.id as user_id, u.username, u.avatar, COUNT(*) as activity_count 
		FROM users u
		LEFT JOIN project_comments pc ON u.id = pc.user_id
		LEFT JOIN reader_comments rc ON u.id = rc.user_id
		LEFT JOIN author_comments ac ON u.id = ac.user_id
		GROUP BY u.id
		ORDER BY activity_count DESC
		LIMIT 5
	`).Scan(&activeUsers)

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"statistics": gin.H{
			"userCount":    userCount,
			"projectCount": projectCount,
			"chapterCount": chapterCount,
			"commentCount": commentCount,
		},
		"recentUsers": recentUsers,
		"hotProjects": hotProjects,
		"activeUsers": activeUsers,
	}))
}

// GetUsers 获取用户列表
func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	role := c.Query("role") // "admin" 或 "user" 或 "全部"

	offset := (page - 1) * pageSize

	var users []pojo.User
	query := config.MysqlDataBase.Model(&pojo.User{})

	// 添加搜索条件
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 根据角色筛选
	if role == "admin" {
		query = query.Where("permission >= ?", 1)
	} else if role == "user" {
		query = query.Where("permission < ?", 1)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	if err := query.Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取用户列表失败"))
		return
	}

	// 隐藏密码
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"users": users,
		"total": total,
	}))
}

// GetUser 获取单个用户信息
func GetUser(c *gin.Context) {
	userID := c.Param("id")

	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "用户不存在"))
		return
	}

	// 隐藏密码
	user.Password = ""

	c.JSON(http.StatusOK, dto.SuccessResponse(user))
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "用户不存在"))
		return
	}

	var updateRequest struct {
		Username   string `json:"username"`
		Email      string `json:"email"`
		Password   string `json:"password"`
		Avatar     string `json:"avatar"`
		Permission uint8  `json:"permission"` // 对应角色：0=普通用户, 1=管理员
		Group      uint8  `json:"group"`      // 对应状态：0=正常, 1=禁用
		Tokens     int    `json:"tokens"`
	}

	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 更新字段
	if updateRequest.Username != "" {
		user.Username = updateRequest.Username
	}
	if updateRequest.Email != "" {
		user.Email = updateRequest.Email
	}
	if updateRequest.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateRequest.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "密码加密失败"))
			return
		}
		user.Password = string(hashedPassword)
	}
	if updateRequest.Avatar != "" {
		user.Avatar = updateRequest.Avatar
	}

	user.Permission = updateRequest.Permission
	user.Group = updateRequest.Group
	user.Tokens = updateRequest.Tokens

	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户信息失败"))
		return
	}

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户信息失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("用户信息更新成功"))
}

// DeleteUser 删除用户 (原接口修改适应POST请求)
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 检查用户是否存在
	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "用户不存在"))
		return
	}

	// 检查是否为管理员
	if user.Permission >= 1 {
		currentUserID, _ := c.Get("userId")
		if currentUserID.(int) != int(user.ID) {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "不能删除其他管理员账户"))
			return
		}
	}

	// 删除用户
	tx := config.MysqlDataBase.Begin()
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除用户失败"))
		return
	}

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除用户失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("用户删除成功"))
}

// UpdateUserStatus 更新用户状态（封禁/解禁）
func UpdateUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var statusRequest struct {
		Status string `json:"status"` // "active" 或 "banned"
	}

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 检查用户是否存在
	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "用户不存在"))
		return
	}

	// 如果尝试禁用管理员账户
	if statusRequest.Status == "banned" && user.Permission >= 1 {
		// 获取当前管理员
		currentUserId, _ := c.Get("userId")
		if currentUserId.(int) != int(user.ID) {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](403, "不能禁用其他管理员账户"))
			return
		}
	}

	// 更新状态 (Group 字段用于表示状态)
	if statusRequest.Status == "active" {
		user.Group = 0
	} else {
		user.Group = 1 // 禁用
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户状态失败"))
		return
	}

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户状态失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("用户状态更新成功"))
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var roleRequest struct {
		Role string `json:"role"` // "user" 或 "admin"
	}

	if err := c.ShouldBindJSON(&roleRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 检查用户是否存在
	var user pojo.User
	if err := config.MysqlDataBase.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "用户不存在"))
		return
	}

	// 更新权限
	if roleRequest.Role == "admin" {
		user.Permission = 1
	} else {
		user.Permission = 0
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户角色失败"))
		return
	}

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新用户角色失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("用户角色更新成功"))
}

// GetChapters 获取所有章节列表（分页、搜索、筛选）
func GetChapters(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	offset := (page - 1) * pageSize

	// 构建查询
	query := config.MysqlDataBase.Model(&pojo.Chapter{}).
		Preload("CurrentVersion").
		Preload("CurrentVersion.User")

	// 按关键词搜索
	if keyword != "" {
		query = query.Where("tittle LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按审核状态过滤
	if status != "" {
		query = query.Joins("JOIN chapter_versions ON chapters.version_id = chapter_versions.id").
			Where("chapter_versions.status = ?", status)
	}

	// 统计总数 - 排除无效版本
	var total int64
	query.Where("version_id IS NOT NULL AND version_id > 0").Count(&total)

	// 获取分页数据 - 排除无效版本
	var chapters []pojo.Chapter
	err := query.Limit(pageSize).Offset(offset).
		Where("version_id IS NOT NULL AND version_id > 0").
		Order("updated_at DESC").
		Find(&chapters).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节列表失败"))
		return
	}

	// 构造返回数据
	result := map[string]interface{}{
		"chapters": chapters,
		"total":    total,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetChapter 获取特定章节详情
func GetChapter(c *gin.Context) {
	chapterID := c.Param("id")

	var chapter pojo.Chapter
	err := config.MysqlDataBase.Preload("CurrentVersion").
		Preload("CurrentVersion.User").
		Where("id = ?", chapterID).
		First(&chapter).Error

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(chapter))
}

// ReviewChapter 审核章节
func ReviewChapter(c *gin.Context) {
	chapterID := c.Param("id")

	var reviewRequest struct {
		Status string `json:"status"` // "approved" 或 "rejected"
		Reason string `json:"reason"` // 拒绝理由（可选）
		Score  int    `json:"score"`  // 内容评分（0-100）
	}

	if err := c.ShouldBindJSON(&reviewRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的请求数据"))
		return
	}

	// 验证评分范围
	if reviewRequest.Score < 0 || reviewRequest.Score > 100 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "评分必须在0-100之间"))
		return
	}

	// 查找章节
	var chapter pojo.Chapter
	err := config.MysqlDataBase.Where("id = ?", chapterID).First(&chapter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	// 查找当前版本
	var currentVersion pojo.ChapterVersion
	err = config.MysqlDataBase.Where("id = ?", chapter.VersionID).First(&currentVersion).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节版本不存在"))
		return
	}

	// 更新版本状态和评分
	currentVersion.Score = reviewRequest.Score
	currentVersion.Status = reviewRequest.Status
	// 如果被拒绝，记录拒绝理由
	if reviewRequest.Status == "rejected" && reviewRequest.Reason != "" {
		// 可以在此处添加拒绝理由字段，或存入日志表
		// 目前模型中没有拒绝理由字段，可以考虑添加
	}

	// 保存更改
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&currentVersion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "审核章节失败"))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "审核章节失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("章节审核成功"))
}

// DeleteChapter 删除章节
func DeleteChapterAdmin(c *gin.Context) {
	chapterID := c.Param("id")

	// 查找章节
	var chapter pojo.Chapter
	err := config.MysqlDataBase.Where("id = ?", chapterID).First(&chapter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 删除关联的版本记录
	if err := tx.Where("chapter_id = ?", chapterID).Delete(&pojo.ChapterVersion{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除章节版本失败"))
		return
	}

	// 删除章节
	if err := tx.Delete(&chapter).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除章节失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除章节失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("章节删除成功"))
}

// UpdateChapterScore 更新章节评分
func UpdateChapterScore(c *gin.Context) {
	chapterID := c.Param("id")

	var scoreRequest struct {
		Score int `json:"score"` // 内容评分（0-100）
	}

	if err := c.ShouldBindJSON(&scoreRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的请求数据"))
		return
	}

	// 验证评分范围
	if scoreRequest.Score < 0 || scoreRequest.Score > 100 {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "评分必须在0-100之间"))
		return
	}

	// 查找章节
	var chapter pojo.Chapter
	err := config.MysqlDataBase.Where("id = ?", chapterID).First(&chapter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	// 查找当前版本
	var currentVersion pojo.ChapterVersion
	err = config.MysqlDataBase.Where("id = ?", chapter.VersionID).First(&currentVersion).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节版本不存在"))
		return
	}

	// 更新评分
	currentVersion.Score = scoreRequest.Score

	// 保存更改
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&currentVersion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新评分失败"))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新评分失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("章节评分更新成功"))
}

// AIScoreChapter 使用AI评估章节内容并给出评分
func AIScoreChapter(c *gin.Context) {
	chapterID := c.Param("id")

	// 查找章节
	var chapter pojo.Chapter
	err := config.MysqlDataBase.Where("id = ?", chapterID).First(&chapter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节不存在"))
		return
	}

	// 查找当前版本
	var currentVersion pojo.ChapterVersion
	err = config.MysqlDataBase.Where("id = ?", chapter.VersionID).First(&currentVersion).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "章节版本不存在"))
		return
	}

	// 获取项目信息为AI提供上下文
	var project pojo.Project
	err = config.MysqlDataBase.Where("id = ?", chapter.ProjectID).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "项目不存在"))
		return
	}

	// 构建提示词
	prompt := fmt.Sprintf(`
请评估以下故事章节的内容，并给出0-100分的评分和详细评分理由。

评分标准：
- 90-100分：优秀，内容健康积极，适合全年龄段
- 70-89分：良好，内容适合青少年，无敏感内容
- 50-69分：一般，包含少量敏感内容，需谨慎
- 0-49分：敏感，包含较多政治、道德、宗教等敏感内容

项目标题：%s
项目类型：%s
章节标题：%s
章节描述：%s

章节内容：
%s

请特别关注内容是否适合青少年，是否符合正确价值观，是否涉及政治、道德、宗教等敏感内容。最后以"评分：X分"和"理由：XXXX"的格式给出你的评估。
`, project.ProjectName, project.Types, chapter.Tittle, chapter.Description, currentVersion.Content)
	var message = []util.Message{
		{
			Role:    "system",
			Content: "你是一个专业的内容审核员，负责评估内容是否适合青少年，是否包含敏感内容。",
		}, {
			Role:    "user",
			Content: prompt,
		},
	}

	// 调用LLM API评估内容 - 参考ProjectService中的调用方式
	maxRetries := 3
	var resp util.ChatResponse

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = util.ChatHandler(util.ChatRequest{
			Model:       util.AgentModelName,
			Messages:    message,
			Question:    prompt,
			Temperature: util.GlobalTemperature,
			MaxTokens:   4000,
		})

		if err == nil {
			break
		}

		// 最后一次尝试失败时返回错误
		if attempt == maxRetries-1 {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "AI评分失败，请稍后重试"))
			return
		}

		// 在重试之前等待一小段时间
		time.Sleep(time.Second * time.Duration(attempt+1))
	}
	res := resp.Choices[0].Message.Content
	// 解析LLM回复，提取评分和理由
	scoreRe := regexp.MustCompile(`评分[:：]\s*(\d+)`)
	reasonRe := regexp.MustCompile(`理由[:：]\s*(.+)(?:\n|$)`)

	scoreMatch := scoreRe.FindStringSubmatch(res)
	reasonMatch := reasonRe.FindStringSubmatch(res)

	var score int = 70 // 默认评分
	var reason string = "AI无法提供明确的评分理由，请人工审核"

	if len(scoreMatch) >= 2 {
		scoreInt, err := strconv.Atoi(scoreMatch[1])
		if err == nil && scoreInt >= 0 && scoreInt <= 100 {
			score = scoreInt
		}
	}

	if len(reasonMatch) >= 2 {
		reason = reasonMatch[1]
	} else {
		// 如果正则匹配失败，使用整个回复作为理由
		reason = res
	}

	// 返回AI评分结果
	result := map[string]interface{}{
		"score":  score,
		"reason": reason,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetProjects 获取所有项目列表（分页、搜索、筛选）
func GetProjects(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")
	typeFilter := c.Query("type")
	statusFilter := c.Query("status") // 添加状态筛选

	offset := (page - 1) * pageSize

	// 构建查询
	query := config.MysqlDataBase.Model(&pojo.Project{}).
		Preload("Team")

	// 按关键词搜索
	if keyword != "" {
		query = query.Where("project_name LIKE ? OR social_story LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 按类型筛选
	if typeFilter != "" {
		query = query.Where("types = ?", typeFilter)
	}

	// 按状态筛选
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}

	// 统计总数
	var total int64
	query.Count(&total)

	// 获取分页数据
	var projects []pojo.Project
	err := query.Limit(pageSize).Offset(offset).Order("updated_at DESC").Find(&projects).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
		return
	}

	// 计算每个项目章节版本的平均评分
	type ProjectScore struct {
		ProjectID  uint
		AvgScore   float64
		ChapterNum int
	}
	var projectScores []ProjectScore

	// 获取所有项目ID
	var projectIDs []uint
	for _, project := range projects {
		projectIDs = append(projectIDs, project.ID)
	}

	// 如果有项目，计算平均分
	if len(projectIDs) > 0 {
		rows, err := config.MysqlDataBase.Raw(`
			SELECT c.project_id, AVG(cv.score) as avg_score, COUNT(DISTINCT c.id) as chapter_num
			FROM chapters c
			JOIN chapter_versions cv ON c.version_id = cv.id
			WHERE c.project_id IN ?
			GROUP BY c.project_id
		`, projectIDs).Rows()

		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var ps ProjectScore
				rows.Scan(&ps.ProjectID, &ps.AvgScore, &ps.ChapterNum)
				projectScores = append(projectScores, ps)
			}
		}
	}

	// 构建项目数据，添加平均分
	type ProjectWithScore struct {
		pojo.Project
		AvgScore   float64 `json:"avg_score"`
		ChapterNum int     `json:"chapter_num"`
	}

	var projectsWithScore []ProjectWithScore
	for _, project := range projects {
		pws := ProjectWithScore{Project: project, AvgScore: 0, ChapterNum: 0}

		// 查找对应项目的得分
		for _, ps := range projectScores {
			if ps.ProjectID == project.ID {
				pws.AvgScore = ps.AvgScore
				pws.ChapterNum = ps.ChapterNum
				break
			}
		}

		projectsWithScore = append(projectsWithScore, pws)
	}

	// 构造返回数据
	result := map[string]interface{}{
		"projects": projectsWithScore,
		"total":    total,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetProject 获取单个项目详情
func GetProject(c *gin.Context) {
	projectID := c.Param("id")

	var project pojo.Project
	err := config.MysqlDataBase.Preload("Team").Where("id = ?", projectID).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "项目不存在"))
		return
	}

	// 获取项目相关的章节
	var chapters []pojo.Chapter
	config.MysqlDataBase.Where("project_id = ?", project.ID).Order("id ASC").Find(&chapters)

	// 获取项目相关的角色
	var characters []pojo.Character
	config.MysqlDataBase.Where("project_id = ?", project.ID).Find(&characters)

	// 计算项目章节的平均评分
	var avgScore float64
	var chapterCount int64
	if len(chapters) > 0 {
		var chapterIDs []uint
		for _, chapter := range chapters {
			if chapter.VersionID > 0 {
				chapterIDs = append(chapterIDs, chapter.ID)
			}
		}

		if len(chapterIDs) > 0 {
			row := config.MysqlDataBase.Raw(`
				SELECT AVG(cv.score) as avg_score, COUNT(DISTINCT c.id) as chapter_count
				FROM chapters c
				JOIN chapter_versions cv ON c.version_id = cv.id
				WHERE c.id IN ?
			`, chapterIDs).Row()

			row.Scan(&avgScore, &chapterCount)
		}
	}

	result := map[string]interface{}{
		"project":       project,
		"chapters":      chapters,
		"characters":    characters,
		"avg_score":     avgScore,
		"chapter_count": chapterCount,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// UpdateProjectStatus 更新项目状态（如推荐显示、禁用等）
func UpdateProjectStatus(c *gin.Context) {
	projectID := c.Param("id")
	var statusRequest struct {
		Status string `json:"status"` // featured, normal, banned
	}

	if err := c.ShouldBindJSON(&statusRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的请求数据"))
		return
	}

	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "项目不存在"))
		return
	}

	// 验证状态值
	validStatus := map[string]bool{
		"featured": true,
		"normal":   true,
		"banned":   true,
	}

	if !validStatus[statusRequest.Status] {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "无效的状态值"))
		return
	}

	// 更新项目状态
	project.Status = statusRequest.Status

	// 保存更改
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新项目状态失败"))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新项目状态失败"))
		return
	}

	// 构建返回消息
	message := "项目状态已更新"
	switch statusRequest.Status {
	case "featured":
		message = "项目已设为推荐"
	case "normal":
		message = "项目已设为正常状态"
	case "banned":
		message = "项目已下架"
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(message))
}

// DeleteProject 删除项目
func DeleteProject(c *gin.Context) {
	projectID := c.Param("id")

	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", projectID).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](404, "项目不存在"))
		return
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 删除相关的章节版本
	// 首先获取所有章节
	var chapters []pojo.Chapter
	if err := tx.Where("project_id = ?", project.ID).Find(&chapters).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	// 获取所有章节的ID
	var chapterIDs []uint
	for _, chapter := range chapters {
		chapterIDs = append(chapterIDs, chapter.ID)
	}

	// 如果有章节，删除相关的章节版本
	if len(chapterIDs) > 0 {
		if err := tx.Where("chapter_id IN ?", chapterIDs).Delete(&pojo.ChapterVersion{}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
			return
		}
	}

	// 删除章节
	if err := tx.Where("project_id = ?", project.ID).Delete(&pojo.Chapter{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	// 删除角色关系
	if err := tx.Exec("DELETE cr FROM character_relation_ships cr JOIN characters c1 ON cr.first_character_id = c1.id WHERE c1.project_id = ?", project.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	// 删除角色
	if err := tx.Where("project_id = ?", project.ID).Delete(&pojo.Character{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	// 删除项目
	if err := tx.Delete(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除项目失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse("项目删除成功"))
}

// GetProjectStats 获取项目统计数据
func GetProjectStats(c *gin.Context) {
	// 项目类型分布
	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}
	config.MysqlDataBase.Model(&pojo.Project{}).
		Select("types as type, COUNT(*) as count").
		Group("types").
		Find(&typeStats)

	// 每月项目创建数量
	var monthlyStats []struct {
		Month string `json:"month"`
		Count int64  `json:"count"`
	}
	config.MysqlDataBase.Model(&pojo.Project{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count").
		Group("month").
		Order("month").
		Find(&monthlyStats)

	// 最活跃的团队（创建项目最多）
	var teamStats []struct {
		TeamID       uint   `json:"team_id"`
		TeamName     string `json:"team_name"`
		ProjectCount int64  `json:"project_count"`
	}
	config.MysqlDataBase.Model(&pojo.Project{}).
		Select("projects.team_id, teams.username as team_name, COUNT(*) as project_count").
		Joins("JOIN teams ON projects.team_id = teams.id").
		Group("projects.team_id").
		Order("project_count DESC").
		Limit(5).
		Find(&teamStats)

	result := map[string]interface{}{
		"typeStats":    typeStats,
		"monthlyStats": monthlyStats,
		"teamStats":    teamStats,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetAdminStatistics 获取管理员数据统计信息
func GetAdminStatistics(c *gin.Context) {
	// 1. 项目相关统计
	var projectStats struct {
		TotalProjects      int64   `json:"total_projects"`
		NormalProjects     int64   `json:"normal_projects"`
		BannedProjects     int64   `json:"banned_projects"`
		AvgScore           float64 `json:"avg_score"`
		TotalWatches       int64   `json:"total_watches"`
		TotalFavorites     int64   `json:"total_favorites"`
		ChapterCount       int64   `json:"chapter_count"`
		HighestRatedID     uint    `json:"highest_rated_id"`
		HighestRatedName   string  `json:"highest_rated_name"`
		HighestRatedScore  float64 `json:"highest_rated_score"`
		MostWatchedID      uint    `json:"most_watched_id"`
		MostWatchedName    string  `json:"most_watched_name"`
		MostWatchedCount   uint    `json:"most_watched_count"`
		MostFavoritedID    uint    `json:"most_favorited_id"`
		MostFavoritedName  string  `json:"most_favorited_name"`
		MostFavoritedCount uint    `json:"most_favorited_count"`
	}

	// 总项目数
	config.MysqlDataBase.Model(&pojo.Project{}).Count(&projectStats.TotalProjects)

	// 正常项目数
	config.MysqlDataBase.Model(&pojo.Project{}).Where("status = ? OR status IS NULL", "normal").Count(&projectStats.NormalProjects)

	// 下架项目数
	config.MysqlDataBase.Model(&pojo.Project{}).Where("status = ?", "banned").Count(&projectStats.BannedProjects)

	// 总观看和收藏数
	config.MysqlDataBase.Model(&pojo.Project{}).Select("SUM(watches) as total_watches, SUM(favorites) as total_favorites").
		Row().Scan(&projectStats.TotalWatches, &projectStats.TotalFavorites)

	// 章节总数
	config.MysqlDataBase.Model(&pojo.Chapter{}).Count(&projectStats.ChapterCount)

	// 计算平均评分
	var totalScore int64
	var ratedChapterCount int64
	rows, err := config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("score > 0").
		Select("SUM(score) as total_score, COUNT(*) as rated_count").
		Rows()

	if err == nil && rows.Next() {
		rows.Scan(&totalScore, &ratedChapterCount)
		if ratedChapterCount > 0 {
			projectStats.AvgScore = float64(totalScore) / float64(ratedChapterCount)
		}
		rows.Close()
	}

	// 评分最高的项目
	var highestRatedProject struct {
		ProjectID uint
		Name      string
		AvgScore  float64
	}

	// 这个查询比较复杂，需要先查出每个章节的最新版本评分，然后按项目分组计算平均分
	subQuery := config.MysqlDataBase.Model(&pojo.Chapter{}).
		Select("project_id, AVG(chapter_versions.score) as avg_score").
		Joins("JOIN chapter_versions ON chapters.version_id = chapter_versions.id").
		Where("chapter_versions.score > 0").
		Group("project_id").
		Order("avg_score DESC").
		Limit(1)

	subQuery.Row().Scan(&highestRatedProject.ProjectID, &highestRatedProject.AvgScore)

	if highestRatedProject.ProjectID > 0 {
		config.MysqlDataBase.Model(&pojo.Project{}).
			Select("project_name").
			Where("id = ?", highestRatedProject.ProjectID).
			Row().Scan(&highestRatedProject.Name)

		projectStats.HighestRatedID = highestRatedProject.ProjectID
		projectStats.HighestRatedName = highestRatedProject.Name
		projectStats.HighestRatedScore = highestRatedProject.AvgScore
	}

	// 观看最多的项目
	var mostWatchedProject pojo.Project
	config.MysqlDataBase.Model(&pojo.Project{}).
		Order("watches DESC").
		First(&mostWatchedProject)

	if mostWatchedProject.ID > 0 {
		projectStats.MostWatchedID = mostWatchedProject.ID
		projectStats.MostWatchedName = mostWatchedProject.ProjectName
		projectStats.MostWatchedCount = mostWatchedProject.Watches
	}

	// 收藏最多的项目
	var mostFavoritedProject pojo.Project
	config.MysqlDataBase.Model(&pojo.Project{}).
		Order("favorites DESC").
		First(&mostFavoritedProject)

	if mostFavoritedProject.ID > 0 {
		projectStats.MostFavoritedID = mostFavoritedProject.ID
		projectStats.MostFavoritedName = mostFavoritedProject.ProjectName
		projectStats.MostFavoritedCount = mostFavoritedProject.Favorites
	}

	// 2. 用户相关统计
	var userStats struct {
		TotalUsers     int64 `json:"total_users"`
		AdminUsers     int64 `json:"admin_users"`
		NormalUsers    int64 `json:"normal_users"`
		BannedUsers    int64 `json:"banned_users"`
		NewUsersToday  int64 `json:"new_users_today"`
		NewUsersWeek   int64 `json:"new_users_week"`
		NewUsersMonth  int64 `json:"new_users_month"`
		MostActiveUser struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Projects int64  `json:"projects"`
		} `json:"most_active_user"`
	}

	// 总用户数
	config.MysqlDataBase.Model(&pojo.User{}).Count(&userStats.TotalUsers)

	// 管理员用户数
	config.MysqlDataBase.Model(&pojo.User{}).Where("permission >= ?", 1).Count(&userStats.AdminUsers)

	// 正常用户数
	config.MysqlDataBase.Model(&pojo.User{}).Where("group = ? OR group IS NULL", 0).Count(&userStats.NormalUsers)

	// 封禁用户数
	config.MysqlDataBase.Model(&pojo.User{}).Where("group = ?", 1).Count(&userStats.BannedUsers)

	// 今日新用户
	config.MysqlDataBase.Model(&pojo.User{}).
		Where("created_at >= ?", time.Now().Format("2006-01-02")).
		Count(&userStats.NewUsersToday)

	// 本周新用户
	weekStart := time.Now().AddDate(0, 0, -int(time.Now().Weekday()))
	config.MysqlDataBase.Model(&pojo.User{}).
		Where("created_at >= ?", weekStart.Format("2006-01-02")).
		Count(&userStats.NewUsersWeek)

	// 本月新用户
	monthStart := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	config.MysqlDataBase.Model(&pojo.User{}).
		Where("created_at >= ?", monthStart.Format("2006-01-02")).
		Count(&userStats.NewUsersMonth)

	// 最活跃用户（创建项目最多）
	config.MysqlDataBase.Raw(`
		SELECT users.id, users.username, COUNT(projects.id) as project_count
		FROM users
		JOIN teams ON users.id = teams.leader_id
		JOIN projects ON teams.id = projects.team_id
		GROUP BY users.id
		ORDER BY project_count DESC
		LIMIT 1
	`).Row().Scan(&userStats.MostActiveUser.ID, &userStats.MostActiveUser.Username, &userStats.MostActiveUser.Projects)

	// 3. 审核相关统计
	var reviewStats struct {
		TotalReviews      int64   `json:"total_reviews"`
		ApprovedReviews   int64   `json:"approved_reviews"`
		RejectedReviews   int64   `json:"rejected_reviews"`
		PendingReviews    int64   `json:"pending_reviews"`
		ApprovalRate      float64 `json:"approval_rate"`
		AvgProcessingTime float64 `json:"avg_processing_time"` // 平均处理时间（天）
		AvgScore          float64 `json:"avg_score"`           // 平均评分
		ReviewsToday      int64   `json:"reviews_today"`
		ReviewsWeek       int64   `json:"reviews_week"`
		ReviewsMonth      int64   `json:"reviews_month"`
	}

	// 总审核数
	//totalVersions := config.MysqlDataBase.Model(&pojo.ChapterVersion{}).Count(&reviewStats.TotalReviews)

	// 通过的审核
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("status = ?", "approved").
		Count(&reviewStats.ApprovedReviews)

	// 拒绝的审核
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("status = ?", "rejected").
		Count(&reviewStats.RejectedReviews)

	// 待审核
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("status = ? OR status IS NULL", "pending").
		Count(&reviewStats.PendingReviews)

	// 通过率
	if reviewStats.TotalReviews > 0 {
		reviewStats.ApprovalRate = float64(reviewStats.ApprovedReviews) / float64(reviewStats.TotalReviews) * 100
	}

	// 平均评分
	if ratedChapterCount > 0 {
		reviewStats.AvgScore = float64(totalScore) / float64(ratedChapterCount)
	}

	// 今日审核数
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("updated_at >= ? AND (status = ? OR status = ?)",
			time.Now().Format("2006-01-02"), "approved", "rejected").
		Count(&reviewStats.ReviewsToday)

	// 本周审核数
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("updated_at >= ? AND (status = ? OR status = ?)",
			weekStart.Format("2006-01-02"), "approved", "rejected").
		Count(&reviewStats.ReviewsWeek)

	// 本月审核数
	config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
		Where("updated_at >= ? AND (status = ? OR status = ?)",
			monthStart.Format("2006-01-02"), "approved", "rejected").
		Count(&reviewStats.ReviewsMonth)

	// 4. 时间趋势数据
	var trendData struct {
		UserRegistration []struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		} `json:"user_registration"`
		ProjectCreation []struct {
			Date  string `json:"date"`
			Count int64  `json:"count"`
		} `json:"project_creation"`
		ChapterReviews []struct {
			Date     string `json:"date"`
			Approved int64  `json:"approved"`
			Rejected int64  `json:"rejected"`
		} `json:"chapter_reviews"`
		ScoreDistribution []struct {
			ScoreRange string `json:"score_range"`
			Count      int64  `json:"count"`
		} `json:"score_distribution"`
	}

	// 用户注册趋势（过去12个月）
	rows, err = config.MysqlDataBase.Raw(`
		SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count
		FROM users
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
		GROUP BY month
		ORDER BY month
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item struct {
				Date  string `json:"date"`
				Count int64  `json:"count"`
			}
			rows.Scan(&item.Date, &item.Count)
			trendData.UserRegistration = append(trendData.UserRegistration, item)
		}
	}

	// 项目创建趋势（过去12个月）
	rows, err = config.MysqlDataBase.Raw(`
		SELECT DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(*) as count
		FROM projects
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
		GROUP BY month
		ORDER BY month
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item struct {
				Date  string `json:"date"`
				Count int64  `json:"count"`
			}
			rows.Scan(&item.Date, &item.Count)
			trendData.ProjectCreation = append(trendData.ProjectCreation, item)
		}
	}

	// 章节审核趋势（过去12个月）
	rows, err = config.MysqlDataBase.Raw(`
		SELECT 
			DATE_FORMAT(updated_at, '%Y-%m') as month, 
			SUM(CASE WHEN status = 'approved' THEN 1 ELSE 0 END) as approved,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected
		FROM chapter_versions
		WHERE updated_at >= DATE_SUB(CURDATE(), INTERVAL 12 MONTH)
		AND (status = 'approved' OR status = 'rejected')
		GROUP BY month
		ORDER BY month
	`).Rows()

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var item struct {
				Date     string `json:"date"`
				Approved int64  `json:"approved"`
				Rejected int64  `json:"rejected"`
			}
			rows.Scan(&item.Date, &item.Approved, &item.Rejected)
			trendData.ChapterReviews = append(trendData.ChapterReviews, item)
		}
	}

	// 评分分布
	scoreRanges := []struct {
		Min   int
		Max   int
		Label string
	}{
		{0, 20, "0-20"},
		{21, 40, "21-40"},
		{41, 60, "41-60"},
		{61, 80, "61-80"},
		{81, 100, "81-100"},
	}

	for _, r := range scoreRanges {
		var count int64
		config.MysqlDataBase.Model(&pojo.ChapterVersion{}).
			Where("score >= ? AND score <= ?", r.Min, r.Max).
			Count(&count)

		trendData.ScoreDistribution = append(trendData.ScoreDistribution, struct {
			ScoreRange string `json:"score_range"`
			Count      int64  `json:"count"`
		}{
			ScoreRange: r.Label,
			Count:      count,
		})
	}

	// 组合所有统计数据
	result := map[string]interface{}{
		"project_stats": projectStats,
		"user_stats":    userStats,
		"review_stats":  reviewStats,
		"trend_data":    trendData,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result))
}

// GetProjectTypeStats 获取项目类型统计
func GetProjectTypeStats(c *gin.Context) {
	var typeStats []struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	config.MysqlDataBase.Model(&pojo.Project{}).
		Select("types as type, COUNT(*) as count").
		Group("types").
		Find(&typeStats)

	c.JSON(http.StatusOK, dto.SuccessResponse(typeStats))
}
