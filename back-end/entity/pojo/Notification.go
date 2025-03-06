/*
 * @Author: Flyinsky w2084151024@gmail.com
 * @Description: None
 */
package pojo

import (
	"time"

	"gorm.io/gorm"
)

// NotificationType 通知类型枚举
type NotificationType int

const (
	SystemNotification        NotificationType = iota + 1 // 系统通知
	LikeNotification                                      // 点赞通知
	CommentNotification                                   // 评论通知
	FollowNotification                                    // 关注通知
	CollaborationNotification                             // 协作邀请通知
	ContentUpdateNotification                             // 内容更新通知
)

// Notification 通知模型
type Notification struct {
	gorm.Model
	Type        NotificationType `gorm:"type:int;not null;comment:通知类型"`  // 通知类型
	Title       string           `gorm:"type:varchar(100);comment:通知标题"`  // 通知标题
	Content     string           `gorm:"type:text;comment:通知内容"`          // 通知内容
	SenderID    uint             `gorm:"comment:发送者ID"`                   // 发送者ID，可为0表示系统通知
	ReceiverID  uint             `gorm:"not null;comment:接收者ID"`          // 接收者ID
	IsRead      bool             `gorm:"default:false;comment:是否已读"`      // 是否已读
	ReadTime    *time.Time       `gorm:"comment:阅读时间"`                    // 阅读时间
	RelatedID   uint             `gorm:"comment:相关内容ID"`                  // 相关内容ID，如作品ID、评论ID等
	RelatedType string           `gorm:"type:varchar(50);comment:相关内容类型"` // 相关内容类型
	ExtraData   string           `gorm:"type:text;comment:额外数据JSON"`      // 额外数据，JSON格式
}

// MessageType 消息类型枚举
type MessageType int

const (
	PrivateMessage MessageType = iota + 1 // 私聊消息
	GroupMessage                          // 群聊消息
)

// Message 私信模型
type Message struct {
	gorm.Model
	Type       MessageType `gorm:"type:int;not null;comment:消息类型"`  // 消息类型
	Content    string      `gorm:"type:text;not null;comment:消息内容"` // 消息内容
	SenderID   uint        `gorm:"not null;comment:发送者ID"`          // 发送者ID
	ReceiverID uint        `gorm:"comment:接收者ID"`                   // 接收者ID，私聊时使用
	GroupID    uint        `gorm:"comment:群组ID"`                    // 群组ID，群聊时使用
	IsRead     bool        `gorm:"default:false;comment:是否已读"`      // 是否已读
	ReadTime   *time.Time  `gorm:"comment:阅读时间"`                    // 阅读时间
	MediaType  string      `gorm:"type:varchar(20);comment:媒体类型"`   // 媒体类型：text, image, audio, video等
	MediaURL   string      `gorm:"type:varchar(255);comment:媒体URL"` // 媒体URL
}

// ChatGroup 聊天群组模型
type ChatGroup struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null;comment:群组名称"` // 群组名称
	Description string `gorm:"type:text;comment:群组描述"`                  // 群组描述
	CreatorID   uint   `gorm:"not null;comment:创建者ID"`                  // 创建者ID
	AvatarURL   string `gorm:"type:varchar(255);comment:群组头像"`          // 群组头像
	MemberCount int    `gorm:"default:1;comment:成员数量"`                  // 成员数量
}

// GroupMember 群组成员模型
type GroupMember struct {
	gorm.Model
	GroupID  uint      `gorm:"not null;comment:群组ID"`               // 群组ID
	UserID   uint      `gorm:"not null;comment:用户ID"`               // 用户ID
	Nickname string    `gorm:"type:varchar(50);comment:群内昵称"`       // 群内昵称
	JoinTime time.Time `gorm:"not null;comment:加入时间"`               // 加入时间
	IsAdmin  bool      `gorm:"default:false;comment:是否为管理员"`        // 是否为管理员
	Status   int       `gorm:"default:1;comment:状态 1:正常 2:禁言 3:退出"` // 成员状态
}

// NotificationSetting 通知设置模型
type NotificationSetting struct {
	gorm.Model
	UserID              uint `gorm:"not null;uniqueIndex;comment:用户ID"` // 用户ID
	SystemNotification  bool `gorm:"default:true;comment:是否接收系统通知"`     // 是否接收系统通知
	LikeNotification    bool `gorm:"default:true;comment:是否接收点赞通知"`     // 是否接收点赞通知
	CommentNotification bool `gorm:"default:true;comment:是否接收评论通知"`     // 是否接收评论通知
	FollowNotification  bool `gorm:"default:true;comment:是否接收关注通知"`     // 是否接收关注通知
	MessageNotification bool `gorm:"default:true;comment:是否接收私信通知"`     // 是否接收私信通知
	EmailNotification   bool `gorm:"default:false;comment:是否接收邮件通知"`    // 是否接收邮件通知
	PushNotification    bool `gorm:"default:true;comment:是否接收推送通知"`     // 是否接收推送通知
}
