package model

import (
	"time"

	"gorm.io/gorm"
)

// name uniqueIndex索引
type User struct {
	Id              int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement;comment:用户id"`
	Name            string `json:"name" gorm:" column:name;uniqueIndex;type:varchar(32);not null;comment:用户名称" valid:"required,stringlength(6|32)"`
	Password        string `json:"password" gorm:"column:password;not null;comment:用户密码"  valid:"required,stringlength(6|32)"`
	Avatar          string `json:"avatar" gorm:"column:avatar;omitempty;comment:用户头像"`
	BackgroundImage string `json:"background_image" gorm:"column:background_image;omitempty;comment:用户个人页顶部大图"`
	Signature       string `json:"signature" gorm:"column:signature;omitempty;comment:个人简介"`

	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt
}

func (User) TableName() string {
	return "user"
}

// 使用 Hook 进行合法性检查
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// if u.FollowCount < 0 || u.FollowerCount < 0 || u.WorkCount < 0 || u.FavoriteCount < 0 || u.TotalFavorited < 0 {
	// 	// logger.Error("Invalid Update: Follow/Unfollow action yield minus value.")
	// 	log.Debug("数据库Count字段为负")
	// 	return fmt.Errorf("操作数为负")
	// }
	return nil
}

//type UserCount struct {
//	Id             int64 `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
//	UserId         int64 `json:"user_id" gorm:"column:user_id;comment:用户id"`
//	FollowCount    int64 `json:"follow_count" gorm:"column:follow_count;omitempty;comment:关注总数"`
//	FollowerCount  int64 `json:"follower_count" gorm:"column:follower_count;omitempty;comment:粉丝总数"`
//	WorkCount      int64 `json:"work_count" gorm:"column:work_count;omitempty;comment:作品数量"`
//	FavoriteCount  int64 `json:"favorite_count" gorm:"column:favorite_count;omitempty;comment:点赞数量"`
//	TotalFavorited int64 `json:"total_favorited" gorm:"column:total_favorited;omitempty;comment:获赞数量"`
//
//	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime"`
//	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime"`
//	DeletedAt gorm.DeletedAt
//}
