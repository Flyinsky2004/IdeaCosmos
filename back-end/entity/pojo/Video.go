package pojo

import "gorm.io/gorm"

//	{
//		"text": "场景文字内容",
//		"illustration_prompt": "插画描述提示词",
//		"image_path": "images/scene1.png",
//		"start_time": 0,
//		"end_time": 15
//	  }
type Scene struct {
	gorm.Model
	Text               string `json:"text"`
	ChapterVersionID   int    `json:"chapter_version_id" bind:"omitempty"`
	IllustrationPrompt string `json:"illustration_prompt"`
	ImagePath          string `json:"image_path"`
	StartTime          int    `json:"start_time"`
	EndTime            int    `json:"end_time"`
}
