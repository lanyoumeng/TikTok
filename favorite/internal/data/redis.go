package data

import (
	"context"
	"errors"
	vpb "favorite/api/video/v1"
	"favorite/internal/pkg/model"
	"favorite/pkg/tool"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

// 获取用户点赞总数
func (f *FavoriteRepo) GetFavoriteCount(ctx context.Context, userId int64) (int64, error) {
	start := time.Now()
	//FavoriteCount userFavList::<userId>的size
	favoriteCount, err := f.data.rdb.SCard(ctx, "userFavList::"+string(userId)).Result()
	if errors.Is(err, redis.Nil) || favoriteCount == 0 {
		//如果缓存不存在，从数据库获取
		err := f.data.db.Model(&model.Favorite{}).Where("user_id = ?", userId).Count(&favoriteCount).Error
		if err != nil {
			f.log.Errorf("GetFavoriteCount-err: %v", err)
			return 0, err
		}

		//缓存
		err = f.data.rdb.Set(context.Background(), "userFavList::"+string(userId), favoriteCount, tool.GetRandomExpireTime()).Err()
		if err != nil {
			f.log.Errorf("GetFavoriteCount-err: %v", err)
			return 0, err

		}

	} else if err != nil {
		f.log.Errorf("GetFavoriteCount-err: %v", err)
		return 0, err
	}

	f.log.Infof("GetFavoriteCount success , GetFavoriteCount耗时=%v", time.Since(start))
	return favoriteCount, nil
}

//作者获赞总数

func (f *FavoriteRepo) GetTotalFavorited(ctx context.Context, authorId int64) (int64, error) {
	start := time.Now()
	var totalFavorited string
	totalFavorited, err := f.data.rdb.Get(ctx, "userTotalFavorited::"+strconv.FormatInt(authorId, 10)).Result()
	if errors.Is(err, redis.Nil) || totalFavorited == "" {
		//如果缓存不存在，从数据库获取
		//获取作者的发布视频id列表
		//汇合所有视频的点赞数
		var videoIds []string
		repo, err := f.data.vc.PublishVidsByAId(context.Background(), &vpb.PublishVidsByAIdReq{AuthorId: authorId})
		if err != nil {
			f.log.Errorf("GetTotalFavorited-err: %v", err)
			return 0, err
		}
		for _, v := range repo.VideoIdList {
			videoIds = append(videoIds, strconv.FormatInt(v, 10))

		}
		var allNum int64
		err = f.data.db.Model(&model.Favorite{}).Where("video_id in (?)", videoIds).Count(&allNum).Error
		if err != nil {
			f.log.Errorf("GetTotalFavorited-err: %v", err)
			return 0, err
		}
		totalFavorited = strconv.FormatInt(allNum, 10)

		//写入缓存
		err = f.data.rdb.Set(context.Background(), "userTotalFavorited::"+strconv.FormatInt(authorId, 10), totalFavorited, tool.GetRandomExpireTime()).Err()
		if err != nil {
			f.log.Errorf("GetTotalFavorited-err: %v", err)
			return 0, err

		}

	} else if err != nil {
		f.log.Errorf("GetTotalFavorited-err: %v", err)
		return 0, err
	}

	f.log.Infof("GetTotalFavorited success , GetTotalFavorited耗时=%v", time.Since(start))
	return strconv.ParseInt(totalFavorited, 10, 64)
}
