package model

import (
	"time"

	"gorm.io/gorm"
)

// id 为索引
type Video struct {
	Id       int64 `json:"id" gorm:"column:id;uniqueIndex;primaryKey;autoIncrement;comment:视频id"`
	AuthorId int64 `json:"author_id" gorm:"column:author_id;not null;comment:作者id"`

	Title    string `json:"title" gorm:"column:title;type:varchar(32);not null;comment:视频标题"`
	PlayUrl  string `json:"play_url" gorm:"column:play_url;not null;comment:视频播放地址"`
	CoverUrl string `json:"cover_url" gorm:"column:cover_url;not null;comment:视频封面地址"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt // 软删除 删除时间
}

func (Video) TableName() string {
	return "video"
}

// 使用 Hook 进行合法性检查
func (u *Video) BeforeUpdate(tx *gorm.DB) error {
	// if u.FollowCount < 0 || u.FollowerCount < 0 || u.WorkCount < 0 || u.FavoriteCount < 0 || u.TotalFavorited < 0 {
	// 	// logger.Error("Invalid Update: Follow/Unfollow action yield minus value.")
	// 	log.Debug("数据库Count字段为负")
	// 	return fmt.Errorf("操作数为负")
	// }
	return nil
}
