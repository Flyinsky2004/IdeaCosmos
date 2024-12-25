package util

import (
	"gopkg.in/gomail.v2"
	"math/rand"
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
