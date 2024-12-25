package pojo

import "gorm.io/gorm"

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:42
 */
type Team struct {
	gorm.Model
	TeamName   string `json:"username;" gorm:"type:varchar(50)"`
	InviteCode string `json:"password;" gorm:"type:varchar(50)"`
	LeaderId   uint   `json:"leaderId;" gorm:"type:bigint"`
}
