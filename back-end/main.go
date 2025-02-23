package main

import (
	"back-end/config"
	"back-end/entity/pojo"
	"back-end/route"
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:10
 */
func main() {
	//printBanner()
	config.InitMysqlDataBase()
	config.InitRedis("localhost:6379", "131598", 0)

	// 禁用外键检查
	//config.MysqlDataBase.Exec("SET FOREIGN_KEY_CHECKS = 0")

	config.MysqlDataBase.AutoMigrate(
		&pojo.User{},
		&pojo.ImageUpload{},
		&pojo.Team{},
		&pojo.JoinRequest{},
		&pojo.Project{},
		&pojo.Character{},
		&pojo.CharacterRelationShip{},
		&pojo.ChapterVersion{},
		&pojo.Chapter{},
		&pojo.ProjectComment{},
		&pojo.AuthorComment{},
		&pojo.ReaderComment{},
		&pojo.Favourite{},
		&pojo.Watch{},
		&pojo.Feeling{},
	)

	// 重新启用外键检查
	//config.MysqlDataBase.Exec("SET FOREIGN_KEY_CHECKS = 1")

	app := gin.Default()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许的前端来源
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-requested-with"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // 预检请求的缓存时间
	}))
	app.Static("/api/uploads", "./uploads")
	app.Static("/api/audio", "./audio")
	route.RegisterRoutes(app)
	app.Run(":8080")
}

func printBanner() {
	// 指定文件路径
	filePath := "resource/banner.txt" // 替换为你的文件路径

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开文件: %v\n", err)
		return
	}
	defer file.Close() // 确保函数结束时关闭文件

	// 使用 bufio.Scanner 逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text()) // 打印每一行内容
	}

	// 检查读取过程中是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Printf("读取文件时发生错误: %v\n", err)
	}
}

// func AiTest() {
// 	//prompt := "A futuristic city with flying cars at sunset"
// 	//baseURL := "https://api1.zhtec.xyz"
// 	//apiKey := "sk-SwmvMY9looEOO7KcEd1a18D8Ad8b413c8c019809586cB842"
// 	//
// 	//imageURL, err := util.GenerateImage(prompt, baseURL, apiKey)
// 	//if err != nil {
// 	//	fmt.Printf("Error generating image: %v\n", err)
// 	//	return
// 	//}
// 	//
// 	//fmt.Printf("Generated image URL: %s\n", imageURL)
// 	//userInfo := util.Message{
// 	//	Role:    "user",
// 	//	Content: "写一个白雪公主大战奥特曼的故事",
// 	//}
// 	message := []util.Message{}
// 	resp, _ := util.ChatHandler(util.ChatRequest{
// 		Model:       "deepseek-chat",
// 		Messages:    message,
// 		Prompt:      "你是一个编剧",
// 		Question:    "写一个白雪公主大战奥特曼的故事",
// 		Temperature: 1.5,
// 	})
// 	jsonStr, _ := json.Marshal(resp)
// 	fmt.Printf(string(jsonStr))
// }
