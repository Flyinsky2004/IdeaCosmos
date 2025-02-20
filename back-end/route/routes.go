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
		publicGroup.GET("getProjectChapters", service.GetProjectChapters)
		publicGroup.GET("getProjectComments", service.GetProjectComments)
		publicGroup.GET("getChapterDetail", service.GetChapterDetail)
		publicGroup.GET("/hot-projects", service.GetHotProjects)
	}

	userGroup := r.Group("/api/user", preHandler())
	{
		userGroup.GET("me", service.GetMyInfo)
		userGroup.POST("updateInfo", service.UpdateUserInfo)
		userGroup.POST("uploadImage", service.UploadImage)
		userGroup.GET("getWebpImageBase64", service.GetImageBase64)
		userGroup.POST("addProjectComment", service.AddProjectComment)
		userGroup.GET("/favorite/add", service.AddFavorite)
		userGroup.GET("/favorite/remove", service.RemoveFavorite)
		userGroup.GET("/favorite/check", service.CheckFavorite)
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
	}

	websocketGroup := r.Group("/api/ws")
	{
		websocketGroup.GET("generateNewChapterVersionStream", service.GenerateNewChapterVersionStream)
	}
	// WebSocket路由
	//r.GET("/ws/chat", service.HandleStreamChat)
}
