package data

import (
	"context"
	"favorite/pkg/tool"
	"github.com/go-redis/redis/v8"
	"strconv"

	// cpb "favorite/api/comment/v1"

	userV1 "favorite/api/user/v1"
	// upb "favorite/api/user/v1"
	vpb "favorite/api/video/v1"
	"favorite/internal/biz"
	"favorite/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"
)

type FavoriteRepo struct {
	data *Data
	log  *log.Helper
}

var _ biz.FavoriteRepo = (*FavoriteRepo)(nil)

// New FavoriteRepo.
func NewBizFavoriteRepo(data *Data, logger log.Logger) biz.FavoriteRepo {
	return &FavoriteRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/favorite")),
	}
}
func (f *FavoriteRepo) GetAuthorIdByVideoId(ctx context.Context, videoId int64) (int64, error) {

	vIds, err := f.data.vc.FavoriteListByVId(ctx, &vpb.FavoriteListReq{VideoIdList: []int64{videoId}})
	if err != nil {
		return 0, err
	}
	return vIds.VideoList[0].Id, nil
}
func (f *FavoriteRepo) FavoriteAction(ctx context.Context, authorId int64, videoId int64, userId int64, actionType int64) error {

	var favorite model.Favorite

	if err := f.data.db.Model(&model.Favorite{}).Where("video_id =? and user_id = ?", videoId, userId).Find(&favorite).Error; err != nil {
		return err
	}

	//执行点赞操作但用户点赞过，直接返回
	if favorite.Id != 0 && actionType == 1 {
		return nil
	}
	//用户未点赞过，执行取消点赞操作，直接返回
	if favorite.Id == 0 && actionType == 2 {
		return nil
	}

	//先更新用户服务缓存 点赞数和被点赞数  防止缓存变化

	//用户的点赞数量
	//FavoriteCount userFavList::<userId>的size
	favoriteCount, err := f.GetFavoriteCount(ctx, userId)
	if err != nil {
		return err
	}
	//作者的被点赞数量
	// TotalFavorited userTotalFavorited::userId
	totalFavorited, err := f.GetTotalFavorited(ctx, authorId)
	if err != nil {
		return err
	}

	favorite.VideoId = videoId
	favorite.UserId = userId

	//1点赞 2取消点赞
	if actionType == 1 {
		//用户的点赞数+1 视频作者的被点赞数+1
		_, err = f.data.userc.UpdateFavoriteCnt(ctx, &userV1.UpdateFavoriteCntRequest{
			FavoriteUserId:  strconv.FormatInt(userId, 10),
			FavoritedUserId: strconv.FormatInt(authorId, 10),
			FavoriteCount:   favoriteCount + 1,
			TotalFavorited:  totalFavorited + 1,
		})
		if err != nil {
			return err
		}

		favorite.Liked = true
		if err := f.data.db.Save(&favorite).Error; err != nil {
			return err
		}

	} else if actionType == 2 {

		//用户的点赞数-1 视频作者的被点赞数-1
		_, err = f.data.userc.UpdateFavoriteCnt(ctx, &userV1.UpdateFavoriteCntRequest{
			FavoriteUserId:  strconv.FormatInt(userId, 10),
			FavoritedUserId: strconv.FormatInt(authorId, 10),
			FavoriteCount:   favoriteCount - 1,
			TotalFavorited:  totalFavorited - 1,
		})
		if err != nil {
			return err
		}

		favorite.Liked = false
		if err := f.data.db.Save(&favorite).Error; err != nil {
			return err
		}

	}

	return nil

}

func (f *FavoriteRepo) FavoriteList(ctx context.Context, userId int64) ([]*vpb.Video, error) {
	// var videos []*vpb.Video

	var videoIds []string
	//缓存获取
	videoIds, err := f.data.rdb.SMembers(ctx, "userFavList::"+strconv.FormatInt(userId, 10)).Result()
	if err != redis.Nil || len(videoIds) == 0 {
		//数据库获取
		err := f.data.db.Model(&model.Favorite{}).Select("video_id").Where("user_id = ?", userId).Find(videoIds).Error
		if err != nil {
			return nil, err
		}

		//写入缓存
		err = f.data.rdb.SAdd(ctx, "userFavList::"+strconv.FormatInt(userId, 10), videoIds).Err()
		if err != nil {
			return nil, err

		}
		//设置过期时间
		err = f.data.rdb.Expire(ctx, "userFavList::"+strconv.FormatInt(userId, 10), tool.GetRandomExpireTime()).Err()
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}
	videoIdList := make([]int64, len(videoIds))
	for i, v := range videoIds {
		videoIdList[i], _ = strconv.ParseInt(v, 10, 64)
	}

	repo, err := f.data.vc.FavoriteListByVId(ctx, &vpb.FavoriteListReq{VideoIdList: videoIdList})
	if err != nil {
		return nil, err
	}

	return repo.VideoList, nil

}

func (f *FavoriteRepo) IsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	var isfavorite string
	//fav::<videoId>::<userId>
	isfavorite, err := f.data.rdb.Get(ctx, "fav::"+strconv.FormatInt(videoId, 10)+"::"+strconv.FormatInt(userId, 10)).Result()
	if err != redis.Nil {
		var favorite model.Favorite
		if err := f.data.db.Model(&model.Favorite{}).Where("video_id =? and user_id = ?", videoId, userId).Find(&favorite).Error; err != nil {
			return false, err
		}
		isfavorite = strconv.FormatBool(favorite.Liked)

		//写入缓存
		err = f.data.rdb.Set(ctx, "fav::"+strconv.FormatInt(videoId, 10)+"::"+strconv.FormatInt(userId, 10), isfavorite, tool.GetRandomExpireTime()).Err()
		if err != nil {
			return false, err

		}

	} else if err != nil {
		return false, err
	}

	return strconv.ParseBool(isfavorite)
}

func (f *FavoriteRepo) GetFavoriteCntByVId(ctx context.Context, videoId int64) (int64, error) {
	var likeCount string
	likeCount, err := f.data.rdb.Get(ctx, "videoFavCnt::"+strconv.FormatInt(videoId, 10)).Result()
	if err != redis.Nil || likeCount == "" {
		var count int64
		err := f.data.db.Model(&model.Favorite{}).Where("video_id=? and liked =?", videoId, true).Count(&count).Error
		if err != nil {
			return 0, err
		}
		likeCount = strconv.FormatInt(count, 10)

		//写入缓存
		err = f.data.rdb.Set(ctx, "videoFavCnt::"+strconv.FormatInt(videoId, 10), likeCount, tool.GetRandomExpireTime()).Err()
		if err != nil {
			return 0, err

		}
	} else if err != nil {
		return 0, err

	}
	return strconv.ParseInt(likeCount, 10, 64)

}

// 获取 用户的 获赞数TotalFavorited 和 点赞数量FavoriteCount
func (f *FavoriteRepo) GetFavoriteCntByUId(ctx context.Context, userId int64) (int64, int64, error) {

	//FavoriteCount userFavList::<userId>的size
	favoriteCount, err := f.GetFavoriteCount(ctx, userId)
	if err != nil {
		return 0, 0, err
	}
	// TotalFavorited userTotalFavorited::userId
	totalFavorited, err := f.GetTotalFavorited(ctx, userId)
	if err != nil {
		return 0, 0, err
	}

	return totalFavorited, favoriteCount, nil
}
