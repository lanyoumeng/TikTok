package model

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
