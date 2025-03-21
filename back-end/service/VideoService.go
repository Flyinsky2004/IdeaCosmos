package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tcolgate/mp3"
)

func getMP3Duration(filePath string) (int, error) {
	// 打开文件
	f, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("无法打开文件: %v", err)
	}
	defer f.Close()

	// 创建mp3解码器
	d := mp3.NewDecoder(f)

	// 计算总时长
	var totalDuration float64
	var frame mp3.Frame
	skipped := 0

	for {
		if err := d.Decode(&frame, &skipped); err != nil {
			if err.Error() == "EOF" { // 文件结束
				break
			}
			return 0, fmt.Errorf("解码错误: %v", err)
		}
		// 获取该帧的时长并累加
		totalDuration += frame.Duration().Seconds()
	}

	// 转换为整数秒并多加1秒
	return int(totalDuration) + 1, nil
}

func GenerateScene(c *gin.Context) {
	chapterVersionID := c.Query("chapter_verison_id")
	var chapterVersion pojo.ChapterVersion
	if err := config.MysqlDataBase.Where("id = ?", chapterVersionID).First(&chapterVersion).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](400, "获取章节版本失败"))
		return
	}
	duration, _ := getMP3Duration("./audio/" + chapterVersion.AudioPath)
	prompt := fmt.Sprintf(`你是一个专业的视频脚本编辑。请将以下故事内容分割成适合制作视频的多个场景。

重要提示：仅返回一个有效的JSON数组，不要包含任何其他文字说明、注释或代码块标记。

输出要求：
1. 每个场景必须包含以下字段：
   - "text": 场景原文内容
   - "illustration_prompt": 详细的插画描述
   - "image_path": 空字符串
   - "start_time": 开始时间（秒）
   - "end_time": 结束时间（秒）
   - "chapter_version_id": 0

音频约束：
- 总音频时长：%d秒
- 使用Azure TTS默认语速每个字约0.3秒
- 所有场景的总时长必须等于音频时长
- 场景切换必须在自然的语句断点处

输入内容：
%s`, duration, chapterVersion.Content)
	question := `严格按照要求只输出JSON数组，不要添加任何额外说明、注释或markdown格式符号。`

	resp, _ := util.ChatHandler(util.ChatRequest{
		Model:       util.AgentModelName,
		Messages:    []util.Message{},
		Prompt:      prompt,
		Question:    question,
		Temperature: util.GlobalTemperature,
		MaxTokens:   4000,
	})

	str := util.CleanJSONResponse(resp.Choices[0].Message.Content)
	var scenes []pojo.Scene

	// 添加重试机制
	maxRetries := 3
	var jsonErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		jsonErr = json.Unmarshal([]byte(str), &scenes)
		if jsonErr == nil {
			break
		}

		// 最后一次尝试失败时返回错误
		if attempt == maxRetries-1 {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "解析场景数据失败，请稍后重试"))
			return
		}

		// 在重试之前等待一小段时间
		time.Sleep(time.Second * time.Duration(attempt+1))
	}

	// 为每个场景设置章节版本ID
	for i := range scenes {
		scenes[i].ChapterVersionID = int(chapterVersion.ID)
	}
	tx := config.MysqlDataBase.Begin()
	// 保存场景到数据库
	if err := tx.Create(&scenes).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "保存场景数据失败，请稍后重试"))
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, dto.SuccessResponse(scenes))
}

func GetSceneByChapterVersionID(c *gin.Context) {
	chapterVersionID := c.Query("chapter_verison_id")
	var scenes []pojo.Scene
	if err := config.MysqlDataBase.Where("chapter_version_id = ?", chapterVersionID).Find(&scenes).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取场景数据失败，请稍后重试"))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse(scenes))
}

func GenerateChapterImages(c *gin.Context) {
	chapterVersionID := c.Query("chapter_verison_id")
	var scenes []pojo.Scene
	if err := config.MysqlDataBase.Where("chapter_version_id = ?", chapterVersionID).Find(&scenes).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取场景数据失败，请稍后重试"))
		return
	}
	tx := config.MysqlDataBase.Begin()
	baseURL := "https://api1.zhtec.xyz"
	apiKey := "sk-SwmvMY9looEOO7KcEd1a18D8Ad8b413c8c019809586cB842"
	for _, scene := range scenes {
		if scene.ImagePath == "" {
			imageURL, err := util.GenerateImage(scene.IllustrationPrompt, baseURL, apiKey)
			if err != nil {
				c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成图片失败，请稍后重试"))
				return
			}
			imagePath, _ := util.DownloadImage(imageURL)
			scene.ImagePath = imagePath
			tx.Save(&scene)
		}
	}
	tx.Commit()
	c.JSON(http.StatusOK, dto.SuccessResponse(scenes))
}

func GenerateChapterVideo(c *gin.Context) {
	chapterVersionID := c.Query("chapter_verison_id")
	var chapterVersion pojo.ChapterVersion
	if err := config.MysqlDataBase.Where("id = ?", chapterVersionID).First(&chapterVersion).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取章节版本失败，请稍后重试"))
		return
	}
	var scenes []pojo.Scene
	if err := config.MysqlDataBase.Where("chapter_version_id = ?", chapterVersionID).Find(&scenes).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取场景数据失败，请稍后重试"))
		return
	}

	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "获取工作目录失败"))
		return
	}

	// 确保视频目录存在
	videoDir := filepath.Join(workDir, "video")
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建视频目录失败"))
		return
	}

	// 生成随机文件名
	randomString := func() string {
		const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
	timestamp := time.Now().Format("20060102")
	videoFilename := fmt.Sprintf("%s_%s.mp4", timestamp, randomString())
	videoPath := filepath.Join(videoDir, videoFilename)

	// 创建临时文件列表
	tempFile, err := os.CreateTemp("", "scenes_*.txt")
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "创建临时文件失败"))
		return
	}
	defer os.Remove(tempFile.Name())

	// 写入场景文件列表
	for _, scene := range scenes {
		// 转换为绝对路径
		imagePath := filepath.Join(workDir, "uploads", scene.ImagePath)
		duration := scene.EndTime - scene.StartTime
		_, err := tempFile.WriteString(fmt.Sprintf("file '%s'\nduration %d\n", imagePath, duration))
		if err != nil {
			c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "写入场景信息失败"))
			return
		}
	}
	tempFile.Close()

	// 获取音频文件绝对路径
	audioPath := filepath.Join(workDir, "audio", chapterVersion.AudioPath)

	// 直接生成带音频的视频
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", tempFile.Name(),
		"-i", audioPath,
		"-vsync", "vfr",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-shortest",
		"-y",
		videoPath)

	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "生成视频失败"+err.Error()))
		return
	}

	// 更新数据库中的视频路径
	if err := config.MysqlDataBase.Model(&chapterVersion).Update("video_path", videoFilename).Error; err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](500, "更新视频路径失败"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(videoFilename))
}
