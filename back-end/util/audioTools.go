package util

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	baseURL      = "https://%s.tts.speech.microsoft.com/cognitiveservices/v1"
	ssmlTemplate = `
<speak version='1.0' xml:lang='%s'>
    <voice xml:lang='%s' xml:gender='%s' name='%s'>
        %s
    </voice>
</speak>`
)

type AzureTTS struct {
	SubscriptionKey string
	Region          string
	client          *http.Client
}

// NewAzureTTS 创建新的Azure TTS客户端
func NewAzureTTS(subscriptionKey, region string) *AzureTTS {
	return &AzureTTS{
		SubscriptionKey: subscriptionKey,
		Region:          region,
		client:          &http.Client{Timeout: 180 * time.Second},
	}
}

// TextToSpeech 将文本转换为语音并保存到文件
func (tts *AzureTTS) TextToSpeech(text, lang, voice, gender string) (string, error) {
	// 生成文件名
	randomString := func() string {
		const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		return string(b)
	}
	timestamp := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s.mp3", timestamp, randomString())
	outputPath := fmt.Sprintf("./audio/%s", filename)

	// 构建SSML
	ssml := fmt.Sprintf(ssmlTemplate, lang, lang, gender, voice, text)

	// 创建请求
	url := fmt.Sprintf(baseURL, tts.Region)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(ssml))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/ssml+xml")
	req.Header.Set("Ocp-Apim-Subscription-Key", tts.SubscriptionKey)
	req.Header.Set("X-Microsoft-OutputFormat", "audio-16khz-128kbitrate-mono-mp3")

	// 发送请求
	resp, err := tts.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API调用失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 确保目录存在
	err = os.MkdirAll("./audio", os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 创建输出文件
	out, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("创建输出文件失败: %v", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("写入文件失败: %v", err)
	}

	return filename, nil
}
