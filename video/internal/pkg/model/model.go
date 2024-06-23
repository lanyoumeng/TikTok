package model

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	Id       int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:视频id"`
	AuthorId int64 `json:"author_id" gorm:"column:author_id;not null;comment:作者id"`

	Title    string `json:"title" gorm:"column:title;type:varchar(32);not null;comment:视频标题"`
	PlayUrl  string `json:"play_url" gorm:"column:play_url;not null;comment:视频播放地址"`
	CoverUrl string `json:"cover_url" gorm:"column:cover_url;not null;comment:视频封面地址"`

	CreatedAt time.Time      `gorm:"column:create_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt // 软删除 删除时间
}

type User struct {
	Id              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:用户id"`
	Name            string `json:"name" gorm:"column:name;unique;type:varchar(32);not null;comment:用户名称"`
	Password        string `json:"password" gorm:"column:password;not null;;comment:用户密码"`
	Avatar          string `json:"avatar" gorm:"column:avatar;omitempty;comment:用户头像"`
	BackgroundImage string `json:"background_image" gorm:"column:background_image;omitempty;comment:用户个人页顶部大图"`
	Signature       string `json:"signature" gorm:"column:signature;omitempty;comment:个人简介"`

	// UserCount UserCount `json:"user_count" gorm:"foreignKey:UserId" `
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

type VideoKafkaMessage struct {
	AuthorId int64  `json:"author_id"` //作者id
	Title    string `json:"title"`     //视频标题

	VideoPath     string `json:"video_path" `      // 视频存储路径
	VideoFileName string `json:"video_file_name" ` // 视频文件名
}

// type VideoList struct {
// 	Video
// 	Author        `json:"author"`
// 	FavoriteCount int64 `json:"favorite_count"`
// 	CommentCount  int64 `json:"comment_count"`
// 	IsFavorite    bool  `json:"is_favorite"`
// }

// type Author struct {
// 	Id              int64  `json:"id"`
// 	Name            string `json:"name"`
// 	FollowCount     int64  `json:"follow_count,omitempty"`
// 	FollowerCount   int64  `json:"follower_count,omitempty"`
// 	IsFollow        bool   `json:"is_follow,omitempty"`
// 	Avatar          string `json:"avatar,omitempty"`
// 	BackgroundImage string `json:"background_image,omitempty"`
// 	Signature       string `json:"signature,omitempty"`
// 	TotalFavorited  int64  `json:"total_favorited,omitempty"`
// 	WorkCount       int64  `json:"work_count,omitempty"`
// 	FavoriteCount   int64  `json:"favorite_count,omitempty"`
// }
