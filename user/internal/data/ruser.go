package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
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
	if err == redis.Nil || len(user) == 0 {
		return nil, errno.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	//刷新过期时间
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		return nil, err
	}
	u := &model.User{}
	err = copier.Copy(u, user)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// 保存用户到redis
func (r *userRepo) RSaveUser(ctx context.Context, user *model.User) error {
	// 保存用户到redis
	userMap, err := tool.StructToMap(user)
	if err != nil {
		return err
	}
	key := "user::" + strconv.FormatInt(user.Id, 10)
	err = r.data.rdb.HSet(ctx, key, userMap).Err()
	if err != nil {
		return err
	}
	// 设置过期时间
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		return err
	}
	return nil

}

//
//func (r *userRepo) RUserInfoAllById(ctx context.Context, id int64) (*v1.User, error) {
//	// 从redis获取所有的用户信息
//	key := "userInfoAll::" + strconv.FormatInt(id, 10)
//	userInfoAll, err := r.data.rdb.HGetAll(ctx, key).Result()
//	if err == redis.Nil || len(userInfoAll) == 0 {
//		return nil, errno.ErrUserNotFound
//	} else if err != nil {
//		return nil, err
//	}
//	//刷新过期时间
//	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
//	if err != nil {
//		return nil, err
//	}
//	userInfo := &v1.User{}
//	err = copier.Copy(userInfo, userInfoAll)
//	if err != nil {
//		return nil, err
//	}
//
//	return userInfo, nil
//
//}
//
//func (r *userRepo) RSaveUserInfoAll(ctx context.Context, user *v1.User) error {
//	//// 保存用户信息到redis
//	//userMap ,err:=structToMap(user)
//	//if err != nil {
//	//	return err
//	//}
//	//err = r.data.rdb.HSet(ctx, strconv.FormatInt(user.Id, 10), userMap).Err()
//	//if err != nil {
//	//	return err
//	//}
//	//// 设置过期时间
//	//_, err = r.data.rdb.Expire(ctx, strconv.FormatInt(user.Id, 10), GetRandomExpireTime()).Result()
//	//if err != nil {
//	//	return err
//	//}
//
//	// Lua 脚本
//	luaScript := `
//    local key = KEYS[1]
//    local expireTime = tonumber(ARGV[1])
//
//    -- 遍历所有的字段和值，并设置到哈希中
//    for i = 2, #ARGV, 2 do
//        local field = ARGV[i]
//        local value = ARGV[i + 1]
//        redis.call('HSET', key, field, value)
//    end
//
//    -- 设置哈希键的过期时间
//    redis.call('EXPIRE', key, expireTime)
//
//    -- 返回结果
//    return 1
//    `
//
//	// 准备 Lua 脚本的参数
//	userKey := "userInfoAll::" + strconv.FormatInt(user.Id, 10) // Redis 键
//	expireTime := int64(tool.GetRandomExpireTime().Seconds())   // 过期时间，以秒为单位
//
//	// 将 User 结构体的字段和值转换为一组参数
//	args := []interface{}{expireTime,
//		"id", user.Id,
//		"name", user.Name,
//		"follow_count", user.FollowCount,
//		"follower_count", user.FollowerCount,
//		"is_follow", strconv.FormatBool(user.IsFollow),
//		"avatar", user.Avatar,
//		"background_image", user.BackgroundImage,
//		"signature", user.Signature,
//		"total_favorited", user.TotalFavorited,
//		"work_count", user.WorkCount,
//		"favorite_count", user.FavoriteCount}
//	// 执行 Lua 脚本
//	_, err := r.data.rdb.Eval(ctx, luaScript, []string{userKey}, args...).Result()
//	if err != nil {
//		fmt.Printf("Error: %v\n", err)
//		return err
//	}
//	return nil
//
//}

func (r *userRepo) RGetCountById(ctx context.Context, id int64) (*model.UserCount, error) {
	// 从redis获取所有计数信息

	key := "count::" + strconv.FormatInt(id, 10)
	count, err := r.data.rdb.HGetAll(ctx, key).Result()
	if err == redis.Nil || len(count) == 0 {
		return nil, errno.ErrUserNotFound
	} else if err != nil {
		return nil, err
	}

	userCount := &model.UserCount{}
	err = copier.Copy(userCount, count)
	if err != nil {
		return nil, err
	}

	return userCount, nil
}
func (r *userRepo) RSaveCount(ctx context.Context, userCount *model.UserCount) error {

	key := "count::" + strconv.FormatInt(userCount.Id, 10)
	countMap, err := tool.StructToMap(userCount)
	if err != nil {
		return err
	}
	err = r.data.rdb.HSet(ctx, key, countMap).Err()
	if err != nil {
		return err
	}
	_, err = r.data.rdb.Expire(ctx, key, tool.GetRandomExpireTime()).Result()
	if err != nil {
		return err
	}
	return nil

}
