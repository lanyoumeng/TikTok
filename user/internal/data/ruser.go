package data

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"user/internal/pkg/errno"
	"user/internal/pkg/model"
	"user/pkg/tool"
)

// 从redis获取用户
func (r *userRepo) RGetUserById(ctx context.Context, id int64) (*model.User, error) {
	// 从redis获取用户
	key := "user::" + strconv.FormatInt(id, 10)
	user, err := r.data.rdb.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) || len(user) == 0 {
		//biz层使用
		return nil, errno.ErrUserNotFound
	} else if err != nil {
		r.log.Error("HGetAll error", err)
		return nil, err
	}
	//刷新过期时间
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		r.log.Error("Expire error", err)
		return nil, err
	}
	u := &model.User{}
	u.Id, _ = strconv.ParseInt(user["id"], 10, 64)
	u.Name = user["name"]
	u.Password = user["password"]
	u.Avatar = user["avatar"]
	u.BackgroundImage = user["background_image"]
	u.Signature = user["signature"]

	//r.log.Infof("RGetUserById: %v", u)
	return u, nil
}

// 保存用户到redis
func (r *userRepo) RSaveUser(ctx context.Context, user *model.User) error {
	// 保存用户到redis
	userMap, err := tool.StructToMap(user)
	if err != nil {
		r.log.Error("StructToMap error", err)
		return err
	}

	key := "user::" + strconv.FormatInt(user.Id, 10)
	err = r.data.rdb.HSet(ctx, key, userMap).Err()
	if err != nil {
		r.log.Error("HSet error", err)
		return err
	}
	// 设置过期时间
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		r.log.Error("Expire error", err)
		return err
	}
	return nil

}

func (r *userRepo) RGetCountById(ctx context.Context, id int64) (*model.UserCount, error) {
	// 从redis获取所有计数信息

	key := "count::" + strconv.FormatInt(id, 10)
	count, err := r.data.rdb.HGetAll(ctx, key).Result()
	if errors.Is(err, redis.Nil) || len(count) == 0 {
		return nil, errno.ErrUserNotFound
	} else if err != nil {
		r.log.Error("HGetAll error", err)
		return nil, err
	}

	userCount := &model.UserCount{}
	//count是被json序列化后存储的，字段名不同 无法直接copy
	//err = copier.Copy(userCount, count)
	//if err != nil {
	//	r.log.Error("copier.Copy error", err)
	//	return nil, err
	//}
	userCount.Id, _ = strconv.ParseInt(count["id"], 10, 64)
	userCount.FollowCount, _ = strconv.ParseInt(count["follow_count"], 10, 64)
	userCount.FollowerCount, _ = strconv.ParseInt(count["follower_count"], 10, 64)
	userCount.WorkCount, _ = strconv.ParseInt(count["work_count"], 10, 64)
	userCount.FavoriteCount, _ = strconv.ParseInt(count["favorite_count"], 10, 64)
	userCount.TotalFavorited, _ = strconv.ParseInt(count["total_favorited"], 10, 64)

	return userCount, nil
}
func (r *userRepo) RSaveCount(ctx context.Context, userCount *model.UserCount) error {

	key := "count::" + strconv.FormatInt(userCount.Id, 10)
	countMap, err := tool.StructToMap(userCount)
	if err != nil {
		r.log.Error("StructToMap error", err)
		return err
	}
	err = r.data.rdb.HSet(ctx, key, countMap).Err()
	if err != nil {
		r.log.Error("HSet error", err)
		return err
	}
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		r.log.Error("Expire error", err)
		return err
	}
	return nil

}
