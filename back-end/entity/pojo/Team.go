package pojo

import "gorm.io/gorm"

/**
 * @author Flyinsky
 * @email w2084151024@gmail.com
 * @date 2024/12/24 09:42
 */
type Team struct {
	gorm.Model
	TeamName        string `json:"username" gorm:"type:varchar(50)"`
	TeamDescription string `json:"teamDescription" gorm:"type:varchar(50)"`
	InviteCode      string `json:"invite_code" gorm:"type:varchar(12)"`
	LeaderId        uint   `json:"leaderId" gorm:"type:bigint"`
}
type JoinRequest struct {
	gorm.Model
	UserId uint `json:"userId" gorm:"not null"`
	TeamId uint `json:"teamId" gorm:"not null"`
	Status int8 `json:"status" gorm:"type:tinyint;default:0"` // 0: pending, 1: approved, 2: rejected
	Team   Team `json:"team" gorm:"foreignKey:TeamId"`
	User   User `json:"user" gorm:"foreignKey:UserId"`
}
type TeamUpdateBody struct {
	ID              uint   `json:"id"`
	TeamName        string `json:"username"`
	TeamDescription string `json:"teamDescription"`
}
type UpdateJoinRequestBody struct {
	RequestId uint `json:"requestId"`
	Status    int8 `json:"status" gorm:"type:tinyint;default:0"` // 0: pending, 1: approved, 2: rejected
}
