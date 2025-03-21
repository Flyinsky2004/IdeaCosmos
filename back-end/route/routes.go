package route

import (
	"back-end/service"

	"github.com/gin-gonic/gin"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:12
 */
func RegisterRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("login", service.Login)
		authGroup.POST("sendCode", service.SendVerifyCode)
		authGroup.POST("register", service.Register)
	}

	publicGroup := r.Group("/api/public")
	{
		publicGroup.GET("getIndexProject", service.GetIndexCoverList)
		publicGroup.GET("getProjectDetail", service.GetProjectDetail)
		publicGroup.GET("getProjectCharacters", service.GetProjectCharacters)
		publicGroup.GET("getProjectCharacterRelationships", service.GetProjectCharacterRelationships)
		publicGroup.GET("getProjectChapters", service.GetProjectChapters)
		publicGroup.GET("getProjectComments", service.GetProjectComments)
		publicGroup.GET("getChapterDetail", service.GetChapterDetail)
		publicGroup.GET("hot-projects", service.GetHotProjects)
		publicGroup.GET("getCategoryProjects", service.GetCategoryProjects)
	}

	userGroup := r.Group("/api/user", preHandler())
	{
		userGroup.GET("me", service.GetMyInfo)
		userGroup.POST("updateInfo", service.UpdateUserInfo)
		userGroup.GET("all", service.GetAllUsers)
		userGroup.POST("uploadImage", service.UploadImage)
		userGroup.GET("getWebpImageBase64", service.GetImageBase64)
		userGroup.POST("addProjectComment", service.AddProjectComment)
		userGroup.POST("addVersionComment", service.AddVersionComment)
		userGroup.GET("getVersionComments", service.GetVersionComments)
		userGroup.GET("/favorite/add", service.AddFavorite)
		userGroup.GET("/favorite/check", service.CheckFavorite)
		userGroup.POST("/feeling/add", service.AddVersionFeeling)
		userGroup.GET("/feeling/get", service.GetVersionFeeling)
		userGroup.GET("/analysis/emotions", service.GetEmotionAnalysis)
	}

	teamGroup := r.Group("/api/team", preHandler())
	{
		teamGroup.POST("createTeam", service.CreateTeam)
		teamGroup.POST("updateTeam", service.UpdateTeam)
		teamGroup.GET("getMyTeams", service.GetMyTeam)
		teamGroup.GET("getMyJoinedTeams", service.GetMyJoinedTeam)
		teamGroup.POST("requestJoin", service.RequestToJoin)
		teamGroup.POST("updateRequest", service.UpdateRequest)
		teamGroup.GET("getRequests", service.GetPendingRequests)
		teamGroup.GET("myTeam", service.GetMyTeam)
		teamGroup.POST("joinByInviteCode", service.JoinTeamByInviteCode)
		teamGroup.GET("members", service.GetTeamMembers)
		teamGroup.GET("regenerateInviteCode", service.RegenerateInviteCode)
		teamGroup.GET("detail", service.GetTeamDetail)
		teamGroup.GET("getTeamByInviteCode", service.GetTeamByInviteCode)
	}

	projectGroup := r.Group("/api/project", preHandler())
	{
		projectGroup.GET("myProjects", service.GetProjectList)
		projectGroup.POST("createProject", service.CreateProject)
		projectGroup.POST("updateProject", service.UpdateProject)
		projectGroup.POST("createCharacter", service.CreateCharacter)
		projectGroup.POST("createCharacterArray", service.CreateCharacterArray)
		projectGroup.POST("updateCharacter", service.UpdateCharacter)
		projectGroup.POST("getCharacters", service.GetCharacters)
		projectGroup.POST("generateCover", service.GenerateProjectCover)
		projectGroup.POST("generateInfo", service.GenerateInfo)
		projectGroup.POST("generateCharacter", service.GenerateCharacter)
		projectGroup.POST("generateCharacterFromDescription", service.GenerateCharacterFromDescription)
		projectGroup.POST("generateCharacterAvatar", service.GenerateCharacterAvatar)
		projectGroup.POST("characterRS/create", service.CreateCharacterRelationship)
		projectGroup.POST("characterRS/update", service.UpdateCharacterRelationship)
		projectGroup.POST("characterRS/delete", service.DeleteCharacterRelationship)
		projectGroup.GET("characterRS/getAll", service.GetCharacterRelationships)
		projectGroup.POST("generateCharacterRS", service.GenerateCharacterRS)
		projectGroup.POST("generateChapters", service.GenerateChapters)
		projectGroup.GET("getAllChapters", service.GetAllChapters)
		projectGroup.POST("createChapterMulti", service.CreateNewChapterMulti)
		projectGroup.POST("generateNewChapterVersion", service.GenerateNewChapterVersion)
		projectGroup.GET("getCurrentChapterVersion", service.GetCurrentChapterVersion)
		projectGroup.POST("createNewChapterVersion", service.CreateNewChapterVersion)
		projectGroup.GET("getChapterVersions", service.GetChapterVersions)
		projectGroup.GET("generateChapterAudio", service.GenerateChapterAudio)
		projectGroup.GET("/analysis/watches-likes", service.GetWatchesAndLikesAnalysis)
		projectGroup.GET("/analysis/style-type", service.GetStyleAndTypeAnalysis)
		projectGroup.POST("creator/comment/add", service.AddCreatorComment)
		projectGroup.GET("creator/comment/list", service.GetCreatorComments)
		projectGroup.POST("creator/comment/delete", service.DeleteCreatorComment)
		projectGroup.POST("chapter/create", service.CreateChapter)
		projectGroup.POST("chapter/update", service.UpdateChapter)
		projectGroup.POST("chapter/delete", service.DeleteChapter)
		projectGroup.POST("deleteCharacter", service.DeleteCharacter)
	}

	// 新增通知服务路由组
	notificationGroup := r.Group("/api/notifications", preHandler())
	{
		notificationGroup.GET("", service.GetUserNotifications)
		notificationGroup.GET("/:id", service.GetNotificationDetail)
		notificationGroup.POST("/:id/read", service.MarkNotificationAsRead)
		notificationGroup.POST("/read-all", service.MarkAllNotificationsAsRead)
		notificationGroup.DELETE("/:id", service.DeleteNotification)
		notificationGroup.DELETE("", service.DeleteAllNotifications)
		notificationGroup.GET("/unread-count", service.GetUnreadNotificationCount)
		notificationGroup.GET("/settings", service.GetNotificationSettings)
		notificationGroup.POST("/settings", service.UpdateNotificationSettings)
	}

	// 新增群组聊天服务路由组
	chatGroup := r.Group("/api/chat", preHandler())
	{
		// 群组管理
		chatGroup.POST("/groups", service.CreateChatGroup)
		chatGroup.GET("/groups", service.GetUserChatGroups)
		chatGroup.GET("/groups/:id", service.GetChatGroupDetail)
		chatGroup.PUT("/groups/:id", service.UpdateChatGroup)
		chatGroup.DELETE("/groups/:id", service.DeleteChatGroup)

		// 成员管理
		chatGroup.POST("/groups/:id/members", service.AddGroupMember)
		chatGroup.DELETE("/groups/:id/members/:userId", service.RemoveGroupMember)
		chatGroup.POST("/groups/:id/leave", service.LeaveGroup)
		chatGroup.GET("/groups/:id/members", service.GetGroupMembers)
		chatGroup.PUT("/groups/:id/members/:userId", service.UpdateGroupMemberInfo)
		chatGroup.POST("/groups/:id/members/:userId/mute", service.MuteGroupMember)
		chatGroup.POST("/groups/:id/members/:userId/admin", service.SetGroupAdmin)

		// 消息管理
		chatGroup.GET("/groups/:id/messages", service.GetGroupMessages)
	}

	// 新增管理员路由组
	adminGroup := r.Group("/api/admin", preHandler(), service.AdminAuthMiddleware())
	{
		adminGroup.GET("/dashboard", service.GetAdminDashboard)

		// 用户管理
		adminGroup.GET("/users", service.GetUsers)
		adminGroup.GET("/users/:id", service.GetUser)
		adminGroup.POST("/users/:id", service.UpdateUser)
		adminGroup.POST("/users/:id/delete", service.DeleteUser)
		adminGroup.POST("/users/:id/status", service.UpdateUserStatus)
		adminGroup.POST("/users/:id/role", service.UpdateUserRole)

		// 章节管理
		adminGroup.GET("/chapters", service.GetChapters)
		adminGroup.GET("/chapters/:id", service.GetChapter)
		adminGroup.POST("/chapters/:id/review", service.ReviewChapter)
		adminGroup.POST("/chapters/:id/delete", service.DeleteChapterAdmin)
		adminGroup.POST("/chapters/:id/score", service.UpdateChapterScore)
		adminGroup.GET("/chapters/:id/ai-score", service.AIScoreChapter)

		// 项目管理
		adminGroup.GET("/projects", service.GetProjects)
		adminGroup.GET("/projects/:id", service.GetProject)
		adminGroup.POST("/projects/:id/status", service.UpdateProjectStatus)
		adminGroup.POST("/projects/:id/delete", service.DeleteProject)
		adminGroup.GET("/projects/stats", service.GetProjectStats)

		// 数据分析
		adminGroup.GET("/statistics/overview", service.GetAdminStatistics)
		adminGroup.GET("/statistics/project-types", service.GetProjectTypeStats)
		adminGroup.POST("/notifications/system", service.SendSystemNotification)
	}

	agentGroup := r.Group("/api/agent", preHandler())
	{
		agentGroup.GET("/chats", service.GetUserChats)
		agentGroup.GET("/chats/:chat_id", service.GetChatHistory)
		agentGroup.GET("/chats/:chat_id/delete", service.DeleteChat)
	}

	websocketGroup := r.Group("/api/ws")
	{
		websocketGroup.GET("generateNewChapterVersionStream", service.GenerateNewChapterVersionStream)
		websocketGroup.GET("modifyChapterVersionStream", service.ModifyChapterVersionStream)
		websocketGroup.GET("newProjectAnalysis", service.NewProjectAnalysis)
		websocketGroup.GET("groupChat/:id", service.HandleGroupChat)
		websocketGroup.GET("projectSuggest", service.IdeaCosmosChat)
	}

	videoGroup := r.Group("/api/video", preHandler())
	{
		videoGroup.GET("generateScene", service.GenerateScene)
		videoGroup.GET("getSceneByChapterVersionID", service.GetSceneByChapterVersionID)
		videoGroup.GET("generateChapterImages", service.GenerateChapterImages)
		videoGroup.GET("generateChapterVideo", service.GenerateChapterVideo)
		videoGroup.GET("getVideosChapterVersion", service.GetVideosChapterVersion)
	}
	// WebSocket路由
	//r.GET("/ws/chat", service.HandleStreamChat)
}
