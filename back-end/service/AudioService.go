/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package service

import (
	"back-end/config"
	"back-end/entity/dto"
	"back-end/entity/pojo"
	"back-end/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var tts = util.NewAzureTTS(
	"384e5a0c65e44b6c8517d149f09438d2",
	"eastasia", // 你的Azure区域
) // 调用TTS服务

func GenerateChapterAudio(c *gin.Context) {
	userId, _ := c.Get("userId")
	chapterId, _ := strconv.ParseInt(c.Query("chapterId"), 10, 64)
	audioName := c.Query("audioName")
	var chapter pojo.Chapter
	err := config.MysqlDataBase.Where("id = ?", chapterId).First(&chapter).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "找不到该篇章"))
		return
	}
	isValidPermission, err := checkProjectPermission(uint(userId.(int)), chapter.ProjectID)
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "校验权限时发生错误"))
		return
	}
	if !isValidPermission {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "您无权访问该项目"))
		return
	}
	var chapterVersion pojo.ChapterVersion
	err = config.MysqlDataBase.Where("id = ?", chapter.VersionID).First(&chapterVersion).Error
	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "没有找到该篇章已有的创作版本"))
		return
	}
	fileName, err := tts.TextToSpeech(
		chapterVersion.Content,
		"zh-CN",   // 语言
		audioName, // 语音名称
		"Female",  // 性别
	)

	if err != nil {
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "生成失败，请稍后重试，错误详细信息:"+err.Error()))
		return
	}
	chapterVersion.AudioPath = fileName
	tx := config.MysqlDataBase.Begin()
	err = tx.Save(&chapterVersion).Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "更新版本时发生错误，错误详细信息:"+err.Error()))
		return
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, dto.ErrorResponse[string](401, "更新版本时发生错误，错误详细信息:"+err.Error()))
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponseWithMessage("生成成功!", fileName))
}
