package tool

import (
	favoritev1 "favorite/api/favorite/v1"
	videov1 "favorite/api/video/v1"
	"math/rand"
	"time"
)

// ConvertVideoList 将 videov1.Video 列表转换为 favoritev1.Video 列表
func ConvertVideoList(videoList []*videov1.Video) []*favoritev1.Video {
	convertedList := make([]*favoritev1.Video, len(videoList))
	for i, video := range videoList {
		convertedList[i] = &favoritev1.Video{
			Id:            video.Id,
			Author:        ConvertUser(video.Author),
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: video.FavoriteCount,
			CommentCount:  video.CommentCount,
			IsFavorite:    video.IsFavorite,
			Title:         video.Title,
		}
	}
	return convertedList
}

// ConvertUser 将 videov1.User 转换为 favoritev1.User
func ConvertUser(user *videov1.User) *favoritev1.User {
	if user == nil {
		return nil
	}
	return &favoritev1.User{
		Id:              user.Id,
		Name:            user.Name,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
}

// redis随机过期时间
func GetRandomExpireTime() time.Duration {
	return time.Duration(300+rand.Intn(300)) * time.Second
}
