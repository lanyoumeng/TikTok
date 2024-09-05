package model

import (
	"time"

	"gorm.io/gorm"
)

// UserId VideoId index索引
type Favorite struct {
	Id      int64 `json:"id" gorm:"primaryKey;autoIncrement"` // 主键
	UserId  int64 `json:"user_id" gorm:"column:user_id;index"`
	VideoId int64 `json:"video_id" gorm:"column:video_id;index"`
	Liked   bool  `json:"liked" gorm:"column:liked"`

	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;autoCreateTime"` // 修改列名
	UpdatedAt time.Time      `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"` // 修改列名
	DeletedAt gorm.DeletedAt // 软删除时间
}

func (Favorite) TableName() string {
	return "favorite"
}

// 使用 Hook 进行合法性检查
func (u *Favorite) BeforeUpdate(tx *gorm.DB) error {
	// if u.FollowCount < 0 || u.FollowerCount < 0 || u.WorkCount < 0 || u.FavoriteCount < 0 || u.TotalFavorited < 0 {
	// 	// logger.Error("Invalid Update: Follow/Unfollow action yield minus value.")
	// 	log.Debug("数据库Count字段为负")
	// 	return fmt.Errorf("操作数为负")
	// }
	return nil
}
