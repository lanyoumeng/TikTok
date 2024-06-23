package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strconv"
	"video/internal/pkg/errno"
	"video/internal/pkg/model"
	"video/pkg/tool"
)

func (v *videoRepo) RZSetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error) {
	var videoId int64
	//ZRevRangeByScoreWithScores(c, key, &redis.ZRangeBy
	//{Min: min, Max: max, Offset: offset, Count: count}).Result()
	videoStr, err := v.data.rdb.ZRevRangeByScoreWithScores(c, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(latestTime, 10),
		Offset: 0,
		Count:  30}).Result()
	if err == redis.Nil || len(videoStr) == 0 {
		return 0, nil, errno.ErrRedisVIdNotFound
	} else if err != nil {
		return 0, nil, err
	}
	videoIds := make([]int64, len(videoStr))
	for i, v := range videoStr {
		videoId, _ = strconv.ParseInt(v.Member.(string), 10, 64)
		videoIds[i] = videoId
	}
	nextTime := int64(videoStr[0].Score)
	return nextTime, videoIds, err
}

// 30个视频id列表存入redis 不设置过期时间
func (v *videoRepo) RZSetSaveVIds(ctx context.Context, value []int64, scores []float64) error {
	//ZAdd(ctx, key, &redis.Z{Score: currentTime, Member: video.Id}).Err()
	z := make([]*redis.Z, len(value))
	for i, val := range value {
		z[i] = &redis.Z{Score: scores[i], Member: strconv.FormatInt(val, 10)}
	}
	err := v.data.rdb.ZAdd(ctx, "videoAll", z...).Err()
	if err != nil {
		return err
	}
	return nil
}

func (v *videoRepo) RGetVideoInfo(ctx context.Context, videoId int64) (*model.Video, error) {
	//HGet(c context.Context, key string, filed string) (string, error)
	key := "videoInfo::" + strconv.FormatInt(videoId, 10)
	video, err := v.data.rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil || len(video) == 0 {
		return nil, errno.ErrRedisVInfoNotFound

	} else if err != nil {
		return nil, err
	}
	videoInfo := &model.Video{}
	videoInfo.Id, _ = strconv.ParseInt(video["id"], 10, 64)
	videoInfo.AuthorId, _ = strconv.ParseInt(video["author_id"], 10, 64)
	videoInfo.Title = video["title"]
	videoInfo.CoverUrl = video["cover_url"]
	videoInfo.PlayUrl = video["play_url"]

	return videoInfo, nil
}
func (v *videoRepo) RSaveVideoInfo(ctx context.Context, video *model.Video) error {
	//HSet(c context.Context, key string, value interface{}) error
	key := "videoInfo::" + strconv.FormatInt(video.Id, 10)
	err := v.data.rdb.HSet(ctx, key, map[string]interface{}{
		"id":        video.Id,
		"author_id": video.AuthorId,
		"title":     video.Title,
		"cover_url": video.CoverUrl,
		"play_url":  video.PlayUrl,
	}).Err()
	if err != nil {
		return err
	}

	//过期时间
	_, err = v.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		return err
	}

	return nil
}

// redis 获取该用户发布视频ID列表
func (v *videoRepo) RPublishVidsByAuthorId(ctx context.Context, authorId int64) ([]int64, error) {
	key := "publishVids::" + strconv.FormatInt(authorId, 10)
	videoIds, err := v.data.rdb.SMembers(ctx, key).Result()
	if err == redis.Nil || len(videoIds) == 0 {
		return nil, errno.ErrRedisPublishVidsNotFound
	} else if err != nil {
		return nil, err
	}

	return tool.StrSliceToInt64Slice(videoIds), nil
}

// redis 该用户发布视频ID列表存入redis
func (v *videoRepo) RSavePublishVids(ctx context.Context, authorId int64, videoIds []int64) error {
	key := "publishVids::" + strconv.FormatInt(authorId, 10)
	strVideoIds := tool.Int64SliceToStrSlice(videoIds)

	// 将 []string 转换为 []interface{}
	interfaces := make([]interface{}, len(strVideoIds))
	for i, v := range strVideoIds {
		interfaces[i] = v
	}

	err := v.data.rdb.SAdd(ctx, key, interfaces...).Err()
	if err != nil {
		return err
	}
	//过期时间
	_, err = v.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		return err
	}
	
	return nil
}
