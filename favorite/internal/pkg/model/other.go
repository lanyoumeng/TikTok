package model

type FavoriteKafkaMessage struct {
	Id      int64 `json:"id" gorm:"primary_key;auto_increment"` // 主键
	UserId  int64 `json:"user_id" gorm:"column:user_id"`
	VideoId int64 `json:"video_id" gorm:"column:video_id"`
	Liked   bool  `json:"liked" gorm:"column:liked"`
}
