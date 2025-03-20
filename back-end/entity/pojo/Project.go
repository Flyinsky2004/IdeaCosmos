package pojo

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/29 16:40
 */
type Project struct {
	gorm.Model
	ProjectName  string         `json:"project_name" gorm:"type:varchar(50)"`
	SocialStory  string         `json:"social_story" gorm:"type:varchar(500)"`
	Start        string         `json:"start" gorm:"type:varchar(400)"`
	HighPoint    string         `json:"high_point" gorm:"type:varchar(400)"`
	Resolved     string         `json:"resolved" gorm:"type:varchar(400)"`
	Style        datatypes.JSON `json:"style" gorm:"type:json"`
	Types        string         `json:"types" gorm:"type:varchar(40)"`
	CoverImage   string         `json:"cover_image" binding:"omitempty" gorm:"type:varchar(200)"`
	MarketPeople datatypes.JSON `json:"market_people" gorm:"type:json"`
	CustomPrompt string         `json:"custom_prompt" gorm:"type:varchar(200)"`
	TeamID       uint           `json:"team_id" gorm:"type:int(11)"`
	Team         Team           `json:"team" bind:"omitempty" gorm:"foreignKey:TeamID"`
	Watches      uint           `json:"watches" gorm:"default:0"`
	Favorites    uint           `json:"favorites" gorm:"default:0"`
	Status       string         `json:"status" gorm:"type:varchar(20);default:'normal'"`
}
type Character struct {
	gorm.Model
	ProjectID   uint    `json:"project_id" gorm:"type:bigint unsigned"`
	Name        string  `json:"name" gorm:"type:varchar(50)"`
	Description string  `json:"description" gorm:"type:varchar(1000)"`
	Avatar      string  `json:"avatar" binding:"omitempty" gorm:"type:varchar(50)"`
	Project     Project `json:"project" gorm:"foreignKey:ProjectID"`
}
type CharacterRelationShip struct {
	gorm.Model
	FirstCharacterID  uint      `json:"first_character_id" gorm:"type:bigint unsigned"`
	SecondCharacterID uint      `json:"second_character_id" gorm:"type:bigint unsigned"`
	FirstCharacter    Character `json:"first_character" gorm:"foreignKey:FirstCharacterID"`
	SecondCharacter   Character `json:"second_character" gorm:"foreignKey:SecondCharacterID"`
	Name              string    `json:"name" gorm:"type:varchar(50)"`
	Content           string    `json:"content" gorm:"type:varchar(1000)"`
}

type Chapter struct {
	gorm.Model
	ProjectID      uint           `json:"project_id" gorm:"type:bigint unsigned"`
	Tittle         string         `json:"Title" gorm:"type:varchar(50)"`
	Description    string         `json:"Description" gorm:"type:varchar(200)"`
	VersionID      uint           `json:"version_id" bind:"omitempty" gorm:"type:bigint unsigned;default:null"`
	CurrentVersion ChapterVersion `json:"current_version" gorm:"foreignKey:VersionID"`
}

type ChapterVersion struct {
	gorm.Model
	UserId           uint   `json:"user_id" gorm:"type:bigint unsigned"`
	ChapterID        uint   `json:"chapter_id" gorm:"type:bigint unsigned"`
	User             User   `json:"user" gorm:"foreignKey:UserId"`
	Content          string `json:"content" gorm:"type:text"`
	OptimizedContent string `json:"optimized_content" gorm:"type:varchar(1000)"`
	AudioPath        string `json:"audio_path" gorm:"type:varchar(20)"`
	VideoPath        string `json:"video_path" gorm:"type:varchar(20)"`
	Score            int    `json:"score" gorm:"type:int(11)"`
	Status           string `json:"status" gorm:"type:varchar(20)"`
}
