package model

// is_followed 是否关注 单独获取
type UserCount struct {
	Id int64 `json:"id" `

	FollowCount   int64 `json:"follow_count" `
	FollowerCount int64 `json:"follower_count"`

	WorkCount int64 `json:"work_count"`

	//用户获赞数
	FavoriteCount int64 `json:"favorite_count"`
	//用户被点赞数
	TotalFavorited int64 `json:"total_favorited"`
}
