package pojo

import "gorm.io/gorm"

type ProjectComment struct {
	gorm.Model
	Content   string `gorm:"type:text"`
	ProjectId uint   `json:"project_id" gorm:"type:bigint unsigned"`
	UserId    uint   `json:"userId" gorm:"type:bigint unsigned"`
	User      User   `json:"user" gorm:"foreignKey:UserId"`
}
type ReaderComment struct {
	gorm.Model
	Content   string `gorm:"type:text"`
	VersionId uint   `json:"version_id" gorm:"type:bigint unsigned"`
	UserId    uint   `json:"userId" gorm:"type:bigint unsigned"`
	User      User   `json:"user" gorm:"foreignKey:UserId"`
}
type AuthorComment struct {
	gorm.Model
	Content   string `gorm:"type:text"`
	VersionId uint   `json:"version_id" gorm:"type:bigint unsigned"`
	UserId    uint   `json:"userId" gorm:"type:bigint unsigned"`
	User      User   `json:"user" gorm:"foreignKey:UserId"`
}
type Favourite struct {
	gorm.Model
	UserId    uint `json:"userId" gorm:"type:bigint unsigned"`
	ProjectId uint `json:"project_id" gorm:"type:bigint unsigned"`
}
type Watch struct {
	gorm.Model
	UserId    uint `json:"userId" gorm:"type:bigint unsigned"`
	ProjectId uint `json:"project_id" gorm:"type:bigint unsigned"`
}
