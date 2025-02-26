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
func DeleteChapter(c *gin.Context) {
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
