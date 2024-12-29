package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:14
 */
var MysqlDataBase *gorm.DB

func InitMysqlDataBase() {
	//dbc := "root:root@tcp(127.0.0.1:9998)/ideacosmos?charset=utf8mb4&parseTime=True&loc=Local"
	dbc := "root:Wjywjy2333@@tcp(127.0.0.1:3306)/ideacosmos?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbc), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库时发生错误:%v", err)
		os.Exit(1)
	}
	MysqlDataBase = db
	fmt.Println("连接数据库成功")
}
