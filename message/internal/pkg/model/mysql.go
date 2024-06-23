package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	Id         int64     `json:"id,omitempty" gorm:"primary_key;autoIncrement:true"`
	ToUserId   int64     `json:"to_user_id,omitempty"`
	FromUserId int64     `json:"from_user_id,omitempty"`
	Content    string    `json:"content,omitempty"`
	CreatedAt  time.Time `json:"create_time" gorm:"column:create_at;autoCreateTime"`

	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (Message) TableName() string {
	return "message"
}

// 使用 Hook 进行合法性检查
func (u *Message) BeforeUpdate(tx *gorm.DB) error {
	// if u.FollowCount < 0 || u.FollowerCount < 0 || u.WorkCount < 0 || u.FavoriteCount < 0 || u.TotalFavorited < 0 {
	// 	// logger.Error("Invalid Update: Follow/Unfollow action yield minus value.")
	// 	log.Debug("数据库Count字段为负")
	// 	return fmt.Errorf("操作数为负")
	// }
	return nil
}
