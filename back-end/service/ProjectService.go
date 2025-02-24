package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetCurrentChapterVersion(c *gin.Context) {
	userId, _ := c.Get("userId")
	chapterId := c.Query("chapter_id")

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", chapterId).First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 检查用户权限
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 如果章节没有关联版本，返回空
	if chapter.VersionID == 0 {
		c.JSON(http.StatusOK, dto.SuccessResponse[string](""))
		return
	}

	// 获取版本内容
	var version pojo.ChapterVersion
	if err := config.MysqlDataBase.Preload("User").Where("id = ?", chapter.VersionID).First(&version).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节版本时发生错误"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(version))
}

// 添加请求消息结构体
type ChapterVersionRequest struct {
	Token      string `json:"token"`
	ChapterId  string `json:"chapterId"`
	WordsCount string `json:"wordsCount"`
}

func GenerateNewChapterVersionStream(c *gin.Context) {
	// 获取基本参数和权限验证
	// userId, _ := c.Get("userId")
	// chapterId := c.Query("chapter_id")
	// wordsCount := c.Query("words_count")
	// 升级HTTP连接到WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "WebSocket升级失败"))
		return
	}
	defer ws.Close()
	// 读取并解析初始消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}
	var request ChapterVersionRequest
	if err := json.Unmarshal(message, &request); err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "无法解析请求参数"+err.Error()))
		return
	}

	// 验证token并获取userId
	claims, err := util.ParseToken(request.Token)
	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		return
	}
	userId := claims.UserID

	chapterId := request.ChapterId
	wordsCount := request.WordsCount
	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", chapterId).First(&chapter).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 验证权限
	isValidPermission, err := checkProjectPermission(uint(userId), chapter.ProjectID)
	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 获取项目信息
	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", chapter.ProjectID).First(&project).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	// 获取角色信息
	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&characters).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取角色信息时发生错误"))
		return
	}

	// 获取角色关系
	var characterIDs []uint
	for _, char := range characters {
		characterIDs = append(characterIDs, char.ID)
	}
	var relationships []pojo.CharacterRelationShip
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取角色关系时发生错误"))
		return
	}

	// 获取所有章节
	var allChapters []pojo.Chapter
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&allChapters).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取章节信息时发生错误"))
		return
	}

	// 构建提示信息
	projectStr := util.ProjectToString(project)
	characterStr := util.CharactersToString(characters)
	characterRelationshipStr := util.CharacterRelationShipsToString(relationships)
	chaptersStr := util.ChaptersToString(allChapters)

	prompt := "【当前章节信息】\n" +
		"标题：" + chapter.Tittle + "\n" +
		"简述：" + chapter.Description + "\n" +
		"要求字数：" + wordsCount + "字\n\n" +
		"【背景参考信息】\n" +
		"项目设定：" + projectStr + "\n" +
		"可用角色：" + characterStr + "\n" +
		"角色关系：" + characterRelationshipStr + "\n" +
		"章节上下文：" + chaptersStr

	systemPrompt := "你是一个专业的" + project.Types + "单章节内容创作者。请你创作当前章节的具体内容。\n\n" +
		"【核心要求】\n" +
		"1. ⚠️ 严格限制：只写当前章节的剧情，不要包含其他章节的具体内容\n" +
		"2. ⚠️ 字数要求：必须严格保证内容（不含标点符号）在" + wordsCount + "字以上\n\n" +
		"【创作指南】\n" +
		"- 把这个章节当作一个独立的短篇，围绕章节简述展开叙述\n" +
		"- 合理使用已有角色，展现他们在本章节中的互动\n" +
		"- 其他章节的内容仅作为背景参考，帮助保持剧情连贯\n" +
		"- 按照项目的风格和设定进行创作\n\n" +
		"【格式要求】\n" +
		"- 直接输出正文内容\n" +
		"- 注意分段，使文章结构清晰\n" +
		"- 不要添加标题、序号或其他额外标记"

	// 调用流式聊天
	streamChan, err := util.StreamChatCompletion(ctx, util.ChatRequest{
		Model:       "deepseek-chat",
		Messages:    []util.Message{},
		Prompt:      systemPrompt,
		Question:    prompt,
		Temperature: 0.6,
		MaxTokens:   8192,
	})

	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 读取流式响应并通过WebSocket发送
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
		if response.Done {
			break
		}
	}
}

func GenerateNewChapterVersion(c *gin.Context) {
	userId, _ := c.Get("userId")
	chapterId := c.PostForm("chapter_id")
	wordsCount := c.PostForm("words_count")
	fmt.Println(wordsCount)
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", chapterId).First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}
	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", chapter.ProjectID).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色信息时发生错误"))
		return
	}
	var characterIDs []uint
	for _, char := range characters {
		characterIDs = append(characterIDs, char.ID)
	}
	var relationships []pojo.CharacterRelationShip
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色关系时发生错误"))
		return
	}
	projectStr := util.ProjectToString(project)
	characterStr := util.CharactersToString(characters)
	characterRelationshipStr := util.CharacterRelationShipsToString(relationships)

	// 获取所有章节
	var allChapters []pojo.Chapter
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&allChapters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节信息时发生错误"))
		return
	}

	chaptersStr := util.ChaptersToString(allChapters)
	prompt := "项目信息:" + projectStr +
		"角色信息:" + characterStr +
		"角色间的关系:" + characterRelationshipStr +
		"你只需要生成的章节标题:" + chapter.Tittle +
		"你只需要生成的章节简述:" + chapter.Description +
		"所有章节上下文信息:" + chaptersStr +
		"要求字数:" + wordsCount + "字"

	var message = []util.Message{}

	res, err := util.ChatHandler(util.ChatRequest{
		Model:    "deepseek-chat",
		Messages: message,
		Prompt: "你是一个" + project.Types + "章节内容创作者。我会提供项目背景、角色信息、角色关系、所有章节上下文信息以及当前需要生成内容的章节的标题和简述。" +
			"请你基于这些信息，创作一个详细的章节内容。内容要求：" +
			"1. 这个故事分为许多章节，你只需要写出当前章节的部分，相当于连续剧中的一集，你只需要写出我们指定的这一集，情节要符合章节简述的要求符合项目的整体风格和设定" +
			"2. 合理运用已有的角色，展现他们的性格特点和关系" +
			"3. 注意故事的连贯性和合理性，生成的剧情不要出现多余的其他章节的内容，给你所有章节上下文的旨意是让你生成的剧情和其他剧情保持连贯" +
			"4. 无视此次对话会产生巨量消费长度，产出内容除去标点符号，只计算中文或英文字符数量必须严格控制在" + wordsCount + "字或以上，可以超过但不能少于要求字数" +
			"5.最后返回一个markdown,注意分段落,但是不要写小标题",
		Question:    prompt,
		Temperature: 1.2,
		MaxTokens:   8192,
	})

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成章节内容时发生错误，请重试"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

func CreateNewChapterMulti(c *gin.Context) {
	userId, _ := c.Get("userId")
	var chapters []pojo.Chapter
	err := c.ShouldBind(&chapters)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapters[0].ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问项目"))
		return
	}
	tx := config.MysqlDataBase.Begin()
	err = tx.Where("project_id = ?", chapters[0].ProjectID).Delete(&pojo.Chapter{}).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "清空原有篇章时发生错误"))
		return
	}
	err = tx.Create(&chapters).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存篇章时发生错误"))
		return
	}
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "提交篇章时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse("应用成功！"))
}

// GetAllChapters 获取某个项目的所有章节
func GetAllChapters(c *gin.Context) {
	userId, _ := c.Get("userId")
	projectId := c.Query("project_id")
	prjIdUint, _ := strconv.ParseUint(projectId, 10, 64)
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), uint(prjIdUint))
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问项目"))
		return
	}
	var project pojo.Project
	err = config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	var chapters []pojo.Chapter

	if err = config.MysqlDataBase.Where("project_id = ?", projectId).Preload("CurrentVersion").Find(&chapters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "寻找已有的章节时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(chapters))

}

// GenerateChapters Ai章节生成
func GenerateChapters(c *gin.Context) {
	userId, _ := c.Get("userId")
	projectId := c.PostForm("project_id")
	prjIdUint, _ := strconv.ParseUint(projectId, 10, 64)
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), uint(prjIdUint))
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问项目"))
		return
	}
	var project pojo.Project
	err = config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	var characters []pojo.Character
	err = config.MysqlDataBase.Where("project_id = ?", projectId).Find(&characters).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "寻找已有的角色时发生错误1"))
		return
	}

	// 查询属于指定项目的角色ID列表
	var characterIDs []uint
	if err := config.MysqlDataBase.Model(&pojo.Character{}).Where("project_id = ?", projectId).Pluck("id", &characterIDs).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "查询角色列表时发生错误"))
		return
	}
	var relationships []pojo.CharacterRelationShip
	// 查询角色关系表中包含这些角色ID的记录
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "查询角色关系时发生错误"))
		return
	}

	projectStr := util.ProjectToString(project)
	characterStr := util.CharactersToString(characters)
	characterRrelationShipStr := util.CharacterRelationShipsToString(relationships)

	prompt := "项目信息:" + projectStr + "角色信息:" + characterStr + "角色间的关系:" + characterRrelationShipStr
	var message = []util.Message{}
	maxRetries := 3
	var res util.ChatResponse
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = util.ChatHandler(util.ChatRequest{
			Model:    "deepseek-chat",
			Messages: message,
			Prompt: "你是一个" + project.Types + "大纲目录设计师，我会提供现有的剧情，角色信息，角色联系等等，你需要基于给出的剧情以及角色背景设计这个作品的章节目录。最后，你需要返回一个json，包含生成的章节目录信息数组,章节目录属性如下，属性名为括号中的英文单词:" +
				"章节标题(Title),章节简述(Description)。其中标题不多于50字，简述不多余200字。",
			Question:    prompt,
			Temperature: 1.3,
			MaxTokens:   8000,
		})

		if err == nil {
			break
		}

		// 最后一次尝试失败时返回错误
		if attempt == maxRetries-1 {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "多次重试后仍然发生错误，请稍后重试"))
			return
		}

		// 在重试之前等待一小段时间
		time.Sleep(time.Second * time.Duration(attempt+1))
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

// GenerateCharacterRS Ai角色关系生成
func GenerateCharacterRS(c *gin.Context) {
	userId, _ := c.Get("userId")
	firstCharacterId := c.PostForm("firstCharacterId")
	secondCharacterId := c.PostForm("secondCharacterId")
	var firstCharacter, secondCharacter *pojo.Character
	err := config.MysqlDataBase.Where("id = ?", firstCharacterId).First(&firstCharacter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "寻找已有的角色时发生错误1"))
		return
	}
	err = config.MysqlDataBase.Where("id = ?", secondCharacterId).First(&secondCharacter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "寻找已有的角色时发生错误2"))
		return
	}
	if firstCharacter.ProjectID != secondCharacter.ProjectID {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "不允许跨项目生成角色关系。"))
		return
	}
	isValidPermisson, err := checkProjectPermission(uint(userId.(int)), firstCharacter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermisson {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问项目"))
		return
	}
	var project pojo.Project
	err = config.MysqlDataBase.Where("id = ?", firstCharacter.ProjectID).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	characterStr := "第一名角色姓名：" + firstCharacter.Name + "其简述为:" + firstCharacter.Description + "第二名角色姓名：" + secondCharacter.Name + "其简述为:" + secondCharacter.Description
	prompt := "受众群体为:" + project.MarketPeople.String() + "两名角色信息为:" + characterStr +
		"内容风格为:" + project.Style.String() + "已有剧情以;隔开：social_story:" + project.SocialStory + ";start" + project.Start + ";high_point" + project.HighPoint + ";resolved" + project.Resolved
	var message = []util.Message{}

	maxRetries := 3
	var res util.ChatResponse

	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = util.ChatHandler(util.ChatRequest{
			Model:    "deepseek-chat",
			Messages: message,
			Prompt: "你是一个" + project.Types + "角色关系设计师，我会提供现有的：社会背景(social_story),开始情景(start),高潮和冲突(high_point)和解决结局(resolved),你需要基于给出的剧情以及角色背景设计两个角色之间的关系。最后，你需要返回一个json，包含生成的角色关系信息,角色关系属性如下，属性名为括号中的英文单词:" +
				"关系名称(name),关系内容(content)，关系名称例如合作伙伴,兄弟,父子,同学等等，关系内容即两名角色之间的故事",
			Question:    prompt,
			Temperature: 1.5,
			MaxTokens:   8000,
		})

		if err == nil {
			break
		}

		// 最后一次尝试失败时返回错误
		if attempt == maxRetries-1 {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "多次重试后仍然发生错误，请稍后重试"))
			return
		}

		// 在重试之前等待一小段时间
		time.Sleep(time.Second * time.Duration(attempt+1))
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

// GenerateCharacterAvatar 生成封面
func GenerateCharacterAvatar(c *gin.Context) {
	characterId := c.PostForm("character_id")
	var character pojo.Character
	err := config.MysqlDataBase.Where("id = ?", characterId).First(&character).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应角色"))
		return
	}
	var project pojo.Project
	err = config.MysqlDataBase.Where("id = ?", character.ProjectID).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	prompt := "生成" + project.Types + "的角色头像，角色名称叫:" + character.Name + "角色介绍:" + character.Description + "这部作品的风格：" + project.Style.String() + "社会背景：" + project.SocialStory + "剧情初始：" + project.Start + "剧情高潮以及核心：" + project.HighPoint + "最后结局：" + project.Resolved
	baseURL := "https://api1.zhtec.xyz"
	apiKey := "sk-SwmvMY9looEOO7KcEd1a18D8Ad8b413c8c019809586cB842"
	imageURL, err := util.GenerateImage(prompt, baseURL, apiKey)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "绘制海报时发生错误，请稍后重试"+"错误信息:"+err.Error()))
		return
	}
	localName, err := util.DownloadImage(imageURL)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存海报时发生错误，请稍后重试 ERR1"+"错误信息:"+err.Error()))
		return
	}
	character.Avatar = localName
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&character).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存海报时发生错误，请稍后重试 ERR2"+"错误信息:"+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存海报时发生错误，请稍后重试 ERR3"+"错误信息:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("角色海报生成成功！", localName))
}

// CreateCharacterArray 批量创建角色
func CreateCharacterArray(c *gin.Context) {
	userId, _ := c.Get("userId")
	var characters []pojo.Character
	if err := c.ShouldBindJSON(&characters); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 设置默认头像
	for i := range characters {
		characters[i].Avatar = "default-avatar.png"
	}

	var projectSource pojo.Project
	if err := config.MysqlDataBase.Where("ID = ?", characters[0].ProjectID).First(&projectSource).Error; err != nil {
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
	if err := tx.Create(&characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存角色信息时发生错误。code : 1"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存角色信息时发生错误。code : 2"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage[string]("角色添加成功！", ""))
}

// GenerateCharacter Ai角色生成
func GenerateCharacter(c *gin.Context) {
	projectId := c.PostForm("project_id")
	var project pojo.Project
	err := config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	var characters []pojo.Character
	err = config.MysqlDataBase.Where("project_id = ?", projectId).Find(&characters).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "寻找已有的角色时发生错误"))
		return
	}
	var characterStr string
	for _, character := range characters {
		characterStr += character.Name + ":" + character.Description + ";"
	}
	prompt := "受众群体为:" + project.MarketPeople.String() + "现有的角色(可能为空，表示没有角色):" + characterStr +
		"内容风格为:" + project.Style.String() + "已有剧情以;隔开：social_story:" + project.SocialStory + ";start" + project.Start + ";high_point" + project.HighPoint + ";resolved" + project.Resolved
	var message = []util.Message{}

	res, err := util.ChatHandler(util.ChatRequest{
		Model:    "deepseek-chat",
		Messages: message,
		Prompt: "你是一个" + project.Types + "角色设计师，我会提供现有的：社会背景(social_story),开始情景(start),高潮和冲突(high_point)和解决结局(resolved),你需要基于给出的剧情设计角色。最后，你需要返回一个json数组，包含生成的所有角色，注意，你生成的结果千万不要包含我给出已有的角色,角色属性如下，属性名为括号中的英文单词:" +
			"姓名(name),描述(description)，对角色的描述包括但不限于性别，人物背景，经历...",
		Question:    prompt,
		Temperature: 1.5,
		MaxTokens:   8000,
	})
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成时发生错误，请重试"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

// GenerateInfo 补全信息
func GenerateInfo(c *gin.Context) {
	projectId := c.PostForm("project_id")
	var project pojo.Project
	err := config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	prompt := "受众群体为:" + project.MarketPeople.String() + "内容风格为:" + project.Style.String() + "已有剧情以;隔开：social_story:" + project.SocialStory + ";start" + project.Start + ";high_point" + project.HighPoint + ";resolved" + project.Resolved
	var message = []util.Message{}

	res, err := util.ChatHandler(util.ChatRequest{
		Model:       "deepseek-chat",
		Messages:    message,
		Prompt:      "你是一个" + project.Types + "补全师，我会提供现有的：社会背景(social_story),开始情景(start),高潮和冲突(high_point)和解决结局(resolved),你需要基于给出的剧情丰富内容，注意这只是故事大概，无需细化，每个属性最多400字。最后，你需要返回一个json,属性名称是括号中的英文单词。",
		Question:    prompt,
		Temperature: 1.5,
		MaxTokens:   8000,
	})
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成时发生错误，请重试"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

// GenerateProjectCover 生成封面
func GenerateProjectCover(c *gin.Context) {
	projectId := c.PostForm("project_id")
	var project pojo.Project
	err := config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}
	prompt := "生成" + project.ProjectName + "的宣传海报，其风格为" + project.Style.String() + "社会背景：" + project.SocialStory + "剧情初始：" + project.Start + "剧情高潮以及核心：" + project.HighPoint + "最后结局：" + project.Resolved
	baseURL := "https://api.gpt.ge"
	apiKey := "sk-hySadfvZfjMxfWx12b302e8c832c4aEeBf7e44C5138bE860"
	imageURL, err := util.GenerateImage(prompt, baseURL, apiKey)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成封面时发生错误，请稍后重试"+"错误信息:"+err.Error()))
		return
	}
	localName, err := util.DownloadImage(imageURL)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存封面时发生错误，请稍后重试 ERR1"+"错误信息:"+err.Error()))
		return
	}
	project.CoverImage = localName
	tx := config.MysqlDataBase.Begin()
	if err := tx.Save(&project).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存封面时发生错误，请稍后重试 ERR2"+"错误信息:"+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存封面时发生错误，请稍后重试 ERR3"+"错误信息:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("封面生成成功！", localName))
}

// GetProjectsByUserId 获取某个用户创建或加入的项目
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

// GetProjectList 获取项目列表
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
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存项目信息时发生错误。code : 1"+err.Error()))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存项目信息时发生错误。code : 2"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse[string]("更新成功"))
}

// GetCharacters 获取角色
func GetCharacters(c *gin.Context) {
	userId, _ := c.Get("userId")
	projectId := c.PostForm("project_id")
	var projectSource pojo.Project
	if err := config.MysqlDataBase.Where("ID = ?", projectId).First(&projectSource).Error; err != nil {
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
	var Characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", projectId).Find(&Characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](501, "查询数据库时发生错误"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(Characters))
}

// CreateCharacter 创建角色
func CreateCharacter(c *gin.Context) {
	userId, _ := c.Get("userId")
	var character pojo.Character
	if err := c.ShouldBindJSON(&character); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	// 设置默认头像
	character.Avatar = "default-avatar.png"

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
func CreateCharacterRelationship(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
		return
	}

	var relationship pojo.CharacterRelationShip
	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}
	// 验证两个角色是否存在并获取项目ID
	var firstChar, secondChar pojo.Character
	if err := config.MysqlDataBase.First(&firstChar, relationship.FirstCharacterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "第一个角色不存在"))
		return
	}
	if err := config.MysqlDataBase.First(&secondChar, relationship.SecondCharacterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "第二个角色不存在"))
		return
	}

	// 验证两个角色是否属于同一个项目
	if firstChar.ProjectID != secondChar.ProjectID {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "两个角色必须属于同一个项目"))
		return
	}

	// 检查权限
	hasPermission, err := checkProjectPermission(uint(userId.(int)), firstChar.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
		return
	}
	if !hasPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限创建角色关系"))
		return
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Create(&relationship).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建角色关系失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建角色关系失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[pojo.CharacterRelationShip](relationship))
}

// UpdateCharacterRelationship 更新角色关系
func UpdateCharacterRelationship(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
		return
	}

	var relationship pojo.CharacterRelationShip
	relationshipId := c.Param("id")

	// 获取原有关系信息
	if err := config.MysqlDataBase.First(&relationship, relationshipId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "角色关系不存在"))
		return
	}

	// 获取角色信息以验证权限
	var character pojo.Character
	if err := config.MysqlDataBase.First(&character, relationship.FirstCharacterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色信息失败"))
		return
	}

	// 检查权限
	hasPermission, err := checkProjectPermission(uint(userId.(int)), character.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
		return
	}
	if !hasPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改角色关系"))
		return
	}

	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Model(&pojo.CharacterRelationShip{}).
		Where("id = ?", relationshipId).
		Updates(relationship).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色关系失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新角色关系失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[string]("更新成功"))
}

// DeleteCharacterRelationship 更新角色关系
func DeleteCharacterRelationship(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
		return
	}

	var relationship pojo.CharacterRelationShip
	relationshipId := c.Param("id")

	// 获取原有关系信息
	if err := config.MysqlDataBase.First(&relationship, relationshipId).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "角色关系不存在"))
		return
	}

	// 获取角色信息以验证权限
	var character pojo.Character
	if err := config.MysqlDataBase.First(&character, relationship.FirstCharacterID).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色信息失败"))
		return
	}

	// 检查权限
	hasPermission, err := checkProjectPermission(uint(userId.(int)), character.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
		return
	}
	if !hasPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限修改角色关系"))
		return
	}

	if err := c.ShouldBindJSON(&relationship); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "请求参数错误"))
		return
	}

	tx := config.MysqlDataBase.Begin()
	if err := tx.Delete(relationship).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除角色关系失败"))
		return
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "删除角色关系失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse[string]("删除成功"))
}

// GetCharacterRelationships 获取角色的所有关系
func GetCharacterRelationships(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "未登录或用户信息获取失败"))
		return
	}

	projectId := c.Query("project_id")
	prjIdUint, _ := strconv.ParseUint(projectId, 10, 64)
	// 检查权限
	hasPermission, err := checkProjectPermission(uint(userId.(int)), uint(prjIdUint))
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证权限时发生错误"))
		return
	}
	if !hasPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您没有权限查看该角色的关系"))
		return
	}

	// 查询属于指定项目的角色ID列表
	var characterIDs []uint
	if err := config.MysqlDataBase.Model(&pojo.Character{}).Where("project_id = ?", projectId).Pluck("id", &characterIDs).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "查询角色列表时发生错误"))
		return
	}
	var relationships []pojo.CharacterRelationShip
	// 查询角色关系表中包含这些角色ID的记录
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "查询角色关系时发生错误"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(relationships))
}

// CreateNewChapterVersion 创建新的章节版本
func CreateNewChapterVersion(c *gin.Context) {
	userId, _ := c.Get("userId")
	var versionRequest struct {
		ChapterID uint   `json:"chapter_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&versionRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", versionRequest.ChapterID).First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 检查用户权限
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 创建新版本
	newVersion := pojo.ChapterVersion{
		UserId:    uint(userId.(int)),
		ChapterID: versionRequest.ChapterID,
		Content:   versionRequest.Content,
	}

	// 开始事务
	tx := config.MysqlDataBase.Begin()

	// 保存新版本
	if err := tx.Create(&newVersion).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存新版本时发生错误"))
		return
	}

	// 更新章节的版本ID
	if err := tx.Model(&chapter).Update("version_id", newVersion.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新章节版本时发生错误"))
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "提交事务时发生错误"))
		return
	}

	// 返回新版本信息
	c.JSON(http.StatusOK, dto.SuccessResponse(newVersion))
}

// OptimizeChapterVersion 根据用户建议优化章节内容
func OptimizeChapterVersion(c *gin.Context) {
	userId, _ := c.Get("userId")
	var optimizeRequest struct {
		ChapterID      uint   `json:"chapter_id" binding:"required"`
		CurrentContent string `json:"current_content" binding:"required"`
		Suggestion     string `json:"suggestion" binding:"required"`
		WordsCount     string `json:"words_count" binding:"required"`
	}

	if err := c.ShouldBindJSON(&optimizeRequest); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "提交了错误的表单"))
		return
	}

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", optimizeRequest.ChapterID).First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 检查用户权限
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 获取项目和角色信息
	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", chapter.ProjectID).First(&project).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&characters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色信息时发生错误"))
		return
	}

	// 获取角色关系
	var characterIDs []uint
	for _, char := range characters {
		characterIDs = append(characterIDs, char.ID)
	}
	var relationships []pojo.CharacterRelationShip
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取角色关系时发生错误"))
		return
	}

	// 获取所有章节
	var allChapters []pojo.Chapter
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&allChapters).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节信息时发生错误"))
		return
	}

	// 构建提示信息
	projectStr := util.ProjectToString(project)
	characterStr := util.CharactersToString(characters)
	characterRelationshipStr := util.CharacterRelationShipsToString(relationships)
	chaptersStr := util.ChaptersToString(allChapters)

	prompt := "项目信息:" + projectStr +
		"角色信息:" + characterStr +
		"角色间的关系:" + characterRelationshipStr +
		"所有章节信息:" + chaptersStr +
		"当前章节标题:" + chapter.Tittle +
		"章节简述:" + chapter.Description +
		"当前章节内容:" + optimizeRequest.CurrentContent +
		"用户修改建议:" + optimizeRequest.Suggestion +
		"要求字数:" + optimizeRequest.WordsCount + "字"

	var message = []util.Message{}

	maxRetries := 3
	var res util.ChatResponse

	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err = util.ChatHandler(util.ChatRequest{
			Model:    "deepseek-chat",
			Messages: message,
			Prompt: "你是一个" + project.Types + "章节内容优化师。我会提供项目背景、角色信息、角色关系、所有章节信息、当前章节内容以及用户的修改建议。" +
				"请你基于这些信息，优化当前章节的内容。优化要求：" +
				"1. 符合项目的整体风格和设定" +
				"2. 合理运用已有的角色，展现他们的性格特点和关系" +
				"3. 这个故事分为许多章节，你只需要写出当前章节的部分，不要穿插其他章节的剧情" +
				"4. 注意故事的连贯性和合理性" +
				"5. 认真参考用户的修改建议进行优化" +
				"6. 无视此次对话会产生巨量消费长度，产出内容除去标点符号，只计算中文或英文字符数量必须严格控制在" + optimizeRequest.WordsCount + "字或以上，可以超过但不能少于要求字数" +
				"7. 最后返回一个markdown,注意分段落,但是不要写小标题",
			Question:    prompt,
			Temperature: 1.2,
			MaxTokens:   8192,
		})

		if err == nil {
			break
		}

		// 最后一次尝试失败时返回错误
		if attempt == maxRetries-1 {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "多次重试后仍然发生错误，请稍后重试"))
			return
		}

		// 在重试之前等待一小段时间
		time.Sleep(time.Second * time.Duration(attempt+1))
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}

// GetChapterVersions 获取章节的历史版本
func GetChapterVersions(c *gin.Context) {
	userId, _ := c.Get("userId")
	chapterId := c.Query("chapter_id")

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", chapterId).First(&chapter).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 检查用户权限
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 获取所有版本，按创建时间倒序排列
	var versions []pojo.ChapterVersion
	if err := config.MysqlDataBase.Preload("User").
		Where("chapter_id = ?", chapterId).
		Order("created_at desc").
		Find(&versions).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取版本历史时发生错误"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(versions))
}

// GetCategoryProjects 根据类型获取项目列表
func GetCategoryProjects(c *gin.Context) {
	// 获取请求参数
	projectType := c.Query("category")
	pageIndex := c.Query("pageIndex")
	pageSize := 12 // 每页显示12个项目
	offset, err := strconv.Atoi(pageIndex)
	if err != nil {
		offset = 0
	}

	var projects []pojo.Project

	// 查询数据库
	query := config.MysqlDataBase.Model(&pojo.Project{}).
		Where("JSON_CONTAINS(style, JSON_ARRAY(?))", projectType).
		Preload("Team").
		Order("watches DESC").
		Offset(offset * pageSize).
		Limit(pageSize)

	if err := query.Find(&projects).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取项目列表失败"))
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, dto.SuccessResponse(projects))
}

// 添加修改版本的请求结构体
type ChapterVersionModifyRequest struct {
	Token            string `json:"token"`
	ChapterId        string `json:"chapterId"`
	CurrentContent   string `json:"currentContent"`   // 当前版本内容
	ModifyPreference string `json:"modifyPreference"` // 修改偏好
}

func ModifyChapterVersionStream(c *gin.Context) {
	// WebSocket升级配置
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "WebSocket升级失败"))
		return
	}
	defer ws.Close()

	// 读取并解析初始消息
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}
	var request ChapterVersionModifyRequest
	if err := json.Unmarshal(message, &request); err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "无法解析请求参数"+err.Error()))
		return
	}

	// 验证token并获取userId
	claims, err := util.ParseToken(request.Token)
	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "token验证失败"))
		return
	}
	userId := claims.UserID

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// 获取章节信息
	var chapter pojo.Chapter
	if err := config.MysqlDataBase.Where("id = ?", request.ChapterId).First(&chapter).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有找到对应章节"))
		return
	}

	// 验证权限
	isValidPermission, err := checkProjectPermission(uint(userId), chapter.ProjectID)
	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "验证用户权限时发生错误"))
		return
	}
	if !isValidPermission {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有权限访问该项目"))
		return
	}

	// 获取项目信息
	var project pojo.Project
	if err := config.MysqlDataBase.Where("id = ?", chapter.ProjectID).First(&project).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	// 获取角色信息
	var characters []pojo.Character
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&characters).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取角色信息时发生错误"))
		return
	}

	// 获取角色关系
	var characterIDs []uint
	for _, char := range characters {
		characterIDs = append(characterIDs, char.ID)
	}
	var relationships []pojo.CharacterRelationShip
	if err := config.MysqlDataBase.Preload("FirstCharacter").Preload("SecondCharacter").
		Where("first_character_id IN ? OR second_character_id IN ?", characterIDs, characterIDs).
		Find(&relationships).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取角色关系时发生错误"))
		return
	}

	// 获取所有章节
	var allChapters []pojo.Chapter
	if err := config.MysqlDataBase.Where("project_id = ?", chapter.ProjectID).Find(&allChapters).Error; err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "获取章节信息时发生错误"))
		return
	}

	// 构建提示信息
	projectStr := util.ProjectToString(project)
	characterStr := util.CharactersToString(characters)
	characterRelationshipStr := util.CharacterRelationShipsToString(relationships)
	chaptersStr := util.ChaptersToString(allChapters)

	prompt := "【❗重要修改要求❗】\n" +
		request.ModifyPreference + "\n\n" +
		"【当前章节信息】\n" +
		"标题：" + chapter.Tittle + "\n" +
		"简述：" + chapter.Description + "\n" +
		"当前版本内容：\n" + request.CurrentContent + "\n\n" +
		"【背景参考信息】\n" +
		"项目设定：" + projectStr + "\n" +
		"可用角色：" + characterStr + "\n" +
		"角色关系：" + characterRelationshipStr + "\n" +
		"章节上下文：" + chaptersStr

	systemPrompt := "你是一个专业的" + project.Types + "内容修改专家。你的首要任务是严格按照用户提供的修改要求进行内容优化。\n\n" +
		"【核心要求】\n" +
		"1. ⚠️ 必须优先满足用户的具体修改要求\n" +
		"2. ⚠️ 在满足修改要求的基础上，保持内容与项目设定的一致性\n" +
		"3. ⚠️ 确保修改符合人物性格和关系的连贯性\n" +
		"4. ⚠️ 修改后的内容字数必须不少于原文字数\n\n" +
		"【修改流程】\n" +
		"1. 首先仔细理解用户的修改要求\n" +
		"2. 严格按照修改要求调整内容\n" +
		"3. 检查修改是否完全符合用户要求\n" +
		"4. 确保修改后字数不少于原文\n" +
		"5. 最后确保与整体项目协调\n\n" +
		"【格式要求】\n" +
		"- 直接输出修改后的正文内容\n" +
		"- 注意分段，使文章结构清晰\n" +
		"- 不要添加标题、序号或其他额外标记"

	// 调用流式聊天
	streamChan, err := util.StreamChatCompletion(ctx, util.ChatRequest{
		Model:       "deepseek-chat",
		Messages:    []util.Message{},
		Prompt:      systemPrompt,
		Question:    prompt,
		Temperature: 0.6,
		MaxTokens:   8192,
	})

	if err != nil {
		ws.WriteJSON(dto.ErrorResponse[string](500, "启动流式生成失败"+err.Error()))
		return
	}

	// 读取流式响应并通过WebSocket发送
	for response := range streamChan {
		if err := ws.WriteJSON(response); err != nil {
			return
		}
		if response.Done {
			break
		}
	}
}

// GenerateCharacterFromDescription 从项目描述中生成角色
func GenerateCharacterFromDescription(c *gin.Context) {
	projectId := c.PostForm("project_id")
	var project pojo.Project
	err := config.MysqlDataBase.Where("id = ?", projectId).First(&project).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "没有找到对应项目"))
		return
	}

	// 获取现有角色，用于避免重复生成
	var existingCharacters []pojo.Character
	err = config.MysqlDataBase.Where("project_id = ?", projectId).Find(&existingCharacters).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取现有角色时发生错误"))
		return
	}

	var existingCharacterStr string
	for _, character := range existingCharacters {
		existingCharacterStr += character.Name + ":" + character.Description + ";"
	}

	// 构建提示信息，重点关注项目描述中的角色信息
	prompt := "【项目基本信息】\n" +
		"项目类型：" + project.Types + "\n" +
		"目标受众：" + project.MarketPeople.String() + "\n" +
		"创作风格：" + project.Style.String() + "\n\n" +
		"【项目剧情信息】\n" +
		"社会背景：" + project.SocialStory + "\n" +
		"开始情节：" + project.Start + "\n" +
		"高潮冲突：" + project.HighPoint + "\n" +
		"结局：" + project.Resolved + "\n\n" +
		"【现有角色】\n" + existingCharacterStr

	var message = []util.Message{}

	res, err := util.ChatHandler(util.ChatRequest{
		Model:    "deepseek-chat",
		Messages: message,
		Prompt: "你是一个专业的" + project.Types + "角色设计师。请仔细分析项目描述中提到的人物，并将其设计为完整的角色。\n\n" +
			"【核心要求】\n" +
			"1. ⚠️ 只生成项目描述中已经明确提到或暗示的角色\n" +
			"2. ⚠️ 不要生成与现有角色重复的角色\n" +
			"3. ⚠️ 为每个角色创建完整的背景故事和性格特征\n\n" +
			"【设计原则】\n" +
			"1. 角色设定要符合项目的整体风格和主题\n" +
			"2. 确保角色背景与项目的社会背景相符\n" +
			"3. 角色性格要能推动剧情发展\n" +
			"4. 注意角色之间的区分度\n\n" +
			"【输出格式】\n" +
			"请返回一个JSON数组，每个角色包含以下属性：\n" +
			"- name: 角色姓名\n" +
			"- description: 角色描述（包含性别、年龄、背景故事、性格特征等）\n" +
			"注意：描述要详细但不超过300字",
		Question:    prompt,
		Temperature: 1.2,
		MaxTokens:   8000,
	})

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成角色时发生错误，请重试"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(res))
}
