package pojo

import "gorm.io/gorm"

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/27 14:13
 */
type ImageUpload struct {
	gorm.Model
	ID     int    `gorm:"primaryKey;autoIncrement;bigint" json:"id"`
	UserId int    `gorm:"bigint" json:"user_id"`
	Path   string `gorm:"varchar(100)" json:"path"`
	Size   int    `gorm:"bigint" json:"size"`
}
