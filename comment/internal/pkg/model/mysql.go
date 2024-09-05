package model

import (
	"gorm.io/gorm"
	"time"
)

// video_id 为Index索引
type Comment struct {
	Id     int64 `json:"id,omitempty" gorm:"column:id;primaryKey;autoIncrement"`
	UserId int64 `json:"user_id" gorm:"column:user_id" `
	//column:password;not null;;comment:用户密码
	VideoId    int64  `json:"video_id,omitempty" gorm:"column:video_id;index"`
	Content    string `json:"content,omitempty" gorm:"column:content"`         //评论内容
	CreateDate string `json:"create_date,omitempty" gorm:"column:create_date"` //评论发布日期，格式 mm-dd

	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (Comment) TableName() string {
	return "comment"
}
