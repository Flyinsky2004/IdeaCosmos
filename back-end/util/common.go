package util

import (
	"back-end/entity/pojo"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 14:59
 */

var senderName = "创剧星球"
var smtpHost = "smtp.office365.com"
var username = "xingongshuo@7jtghd.onmicrosoft.com"
var password = "Wjywjy2333%^"
var smtpPort = 587

// GenerateCode 生成随机验证码
func GenerateCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", senderName)
	m.SetHeader("From", username)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, username, password)

	return d.DialAndSend(m)
}
func GenerateRandomString(n int) string {
	// 定义随机字符集合
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// 初始化随机种子
	rand.Seed(time.Now().UnixNano())

	// 使用strings.Builder优化字符串拼接
	var sb strings.Builder
	sb.Grow(n)
	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}
	return sb.String()
}

// ProjectToString 将 Project 转为字符串，包含所有信息
func ProjectToString(project pojo.Project) string {
	return fmt.Sprintf(
		`%s 名称: %s; 
社会背景: %s; 
剧情初始: %s; 
冲突/高光时刻: %s; 
解决结局: %s; 
风格: %s; 
受众群体: %s;`,
		project.Types,
		project.ProjectName,
		project.SocialStory,
		project.Start,
		project.HighPoint,
		project.Resolved,
		string(project.Style), // 风格 (JSON 格式)
		string(project.MarketPeople))
}

// CharacterToString 将 Character 转为字符串
func CharacterToString(character pojo.Character) string {
	return fmt.Sprintf(
		"角色名称: %s; 描述: %s;",
		character.Name, character.Description,
	)
}

// CharacterRelationShipToString 将 CharacterRelationShip 转为字符串
func CharacterRelationShipToString(relation pojo.CharacterRelationShip) string {
	return fmt.Sprintf(
		"%s和%s的关系是: %s; 详细信息: %s;",
		relation.FirstCharacter.Name, relation.SecondCharacter.Name, relation.Name, relation.Content,
	)
}

// CharacterRelationShipsToString 将 CharacterRelationShip 数组转为字符串
func CharacterRelationShipsToString(relations []pojo.CharacterRelationShip) string {
	var result string
	for _, relation := range relations {
		result += fmt.Sprintf(
			"%s和%s的关系是: %s; 详细信息: %s;\n",
			relation.FirstCharacter.Name, relation.SecondCharacter.Name, relation.Name, relation.Content,
		)
	}
	return result
}

// CharactersToString 将 Character 数组转为字符串
func CharactersToString(characters []pojo.Character) string {
	var result string
	for _, character := range characters {
		result += fmt.Sprintf(
			"角色名称: %s; 描述: %s;\n",
			character.Name, character.Description,
		)
	}
	return result
}

// ChaptersToString 将 Chapter 数组转为字符串
func ChaptersToString(chapters []pojo.Chapter) string {
	var result string
	for i, chapter := range chapters {
		result += fmt.Sprintf(
			"第%d章 标题: %s; 简述: %s;\n",
			i+1, chapter.Tittle, chapter.Description,
		)
	}
	return result
}
