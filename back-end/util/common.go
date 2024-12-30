package util

import (
	"gopkg.in/gomail.v2"
	"math/rand"
	"strings"
	"time"
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
