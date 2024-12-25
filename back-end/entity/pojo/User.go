package pojo

import "gorm.io/gorm"

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:40
 */
type User struct {
	gorm.Model
	Username   string `json:"username;" gorm:"type:varchar(50)"`
	Password   string `json:"password;" gorm:"type:varchar(50)"`
	Email      string `json:"email;" gorm:"type:varchar(50)"`
	Tokens     int    `json:"tokens;" gorm:"type:MEDIUMINT"`
	Permission uint8  `json:"permission;" gorm:"type:tinyint"`
	Group      uint8  `json:"group;" gorm:"type:tinyint"`
}
