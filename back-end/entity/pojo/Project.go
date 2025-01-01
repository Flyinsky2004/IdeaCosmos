package pojo

import "gorm.io/gorm"

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/29 16:40
 */
type Project struct {
	gorm.Model
	ProjectName  string   `json:"project_name" gorm:"type:varchar(50)"`
	SocialStory  string   `json:"social_story" gorm:"type:varchar(500)"`
	Start        string   `json:"start" gorm:"type:varchar(400)"`
	HighPoint    string   `json:"high_point" gorm:"type:varchar(400)"`
	Resolved     string   `json:"resolved" gorm:"type:varchar(400)"`
	Style        []string `json:"style" gorm:"type:json"`
	Types        string   `json:"types" gorm:"type:varchar(40)"`
	CoverImage   string   `json:"cover_image" binding:"omitempty" gorm:"type:varchar(200)"`
	MarketPeople []string `json:"market_people" gorm:"type:json"`
	CustomPrompt string   `json:"custom_prompt" gorm:"type:varchar(200)"`
	TeamID       uint     `json:"team_id" gorm:"type:int(11)"`
	Team         Team     `json:"team" gorm:"foreignKey:TeamID"`
}
type Character struct {
	gorm.Model
	ProjectID   uint    `json:"project_id" gorm:"type:bigint unsigned"`
	Name        string  `json:"role_name" gorm:"type:varchar(50)"`
	Description string  `json:"description" gorm:"type:varchar(200)"`
	Project     Project `json:"project" gorm:"foreignKey:ProjectID"`
}
type CharacterRelationShip struct {
	gorm.Model
	FirstCharacterID  uint      `json:"first_character_id" gorm:"type:bigint unsigned"`
	SecondCharacterID uint      `json:"second_character_id" gorm:"type:bigint unsigned"`
	FirstCharacter    Character `json:"first_character" gorm:"foreignKey:FirstCharacterID"`
	SecondCharacter   Character `json:"second_character" gorm:"foreignKey:SecondCharacterID"`
	Content           string    `json:"content" gorm:"type:varchar(200)"`
}
