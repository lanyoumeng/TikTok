package model

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
	Id         int64  `json:"id,omitempty" gorm:"primaryKey;autoIncrement:true"`
	UserId     int64  `json:"user_id"`
	VideoId    int64  `json:"video_id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"` //评论发布日期，格式 mm-dd

	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (Comment) TableName() string {
	return "comment"
}
