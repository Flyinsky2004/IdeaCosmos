package main

import (
	"back-end/config"
	"back-end/entity/pojo"
	"back-end/route"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:10
 */
func main() {
	printBanner()
	config.InitMysqlDataBase()
	config.InitRedis("localhost:9999", "131598", 0)
	config.MysqlDataBase.AutoMigrate(&pojo.User{})
	app := gin.Default()
	app.Use(route.CorsHandler())
	app.Use(gin.Logger())
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
