package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	messageV1 "relation/api/message/v1"
	pb "relation/api/relation/v1"
	userV1 "relation/api/user/v1"
	"relation/internal/biz"
	"relation/internal/pkg/model"
	"strconv"
)

type relationRepo struct {
	data *Data
	log  *log.Helper
}

// NewRelationRepo .
func NewRelationRepo(data *Data, logger log.Logger) biz.RelationRepo {
	return &relationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *relationRepo) IsFollow(ctx context.Context, userId, targetUserId int64) (bool, error) {
	//set  follow::user_id  关注用户的ids
	//1.缓存获取
	flag, err := r.data.rdb.SIsMember(ctx, "follow::"+strconv.Itoa(int(userId)), targetUserId).Result()
	if err == redis.Nil {
		//2.数据库获取
		var count int64
		r.data.db.Model(&model.Relation{}).Where("user_id = ? AND to_user_id = ?", userId, targetUserId).Count(&count)
		if count > 0 {
			// 记录存在
			flag = true

			//更新缓存
			err := r.data.rdb.SAdd(ctx, "follow::"+strconv.Itoa(int(userId)), targetUserId).Err()
			if err != nil {
				return false, err
			}

		} else {
			// 记录不存在
			flag = false

			//更新缓存
			err := r.data.rdb.SRem(ctx, "follow::"+strconv.Itoa(int(userId)), targetUserId).Err()
			if err != nil {
				return false, err
			}

		}
	} else if err != nil {
		return false, err
	}
	return flag, nil

}
func (r *relationRepo) InsertFollow(ctx context.Context, userId, targetUserId int64) error {
	//1.插入数据库
	err := r.data.db.Create(&model.Relation{UserId: userId, ToUserId: targetUserId}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *relationRepo) DeleteFollow(ctx context.Context, userId, targetUserId int64) error {
	//1.删除数据库
	err := r.data.db.Where("user_id = ? AND to_user_id = ?", userId, targetUserId).Delete(&model.Relation{}).Error
	if err != nil {
		return err
	}

	return nil
}
func (r *relationRepo) FollowUserIdList(ctx context.Context, userId int64) ([]int64, error) {
	//1.缓存获取
	var userIds []int64

	ids, err := r.data.rdb.SMembers(ctx, "follow::"+strconv.Itoa(int(userId))).Result()
	if err == redis.Nil {
		//2.数据库获取
		r.data.db.Model(&model.Relation{}).Select("user_id").Where("user_id = ?", userId).Find(&userIds)

		//更新缓存
		err := r.data.rdb.SAdd(ctx, "follow::"+strconv.Itoa(int(userId)), userIds).Err()
		if err != nil {
			return nil, err
		}

		return userIds, err

	} else if err != nil {
		return nil, err
	}

	for _, id := range ids {
		userId, _ := strconv.ParseInt(id, 10, 64)
		userIds = append(userIds, userId)
	}

	return userIds, nil
}

func (r *relationRepo) UserInfoList(ctx context.Context, userIdList []int64) ([]*pb.User, error) {

	//获取关注用户的ids 然后rpc获取用户信息
	userInfoList, err := r.data.userc.UserInfoList(ctx, &userV1.UserInfoListrRequest{UserId: userIdList})
	if err != nil {
		return nil, err
	}

	pbUserInfoList := make([]*pb.User, len(userInfoList.Users))
	for i, userInfo := range userInfoList.Users {
		err = copier.Copy(pbUserInfoList[i], userInfo)
		if err != nil {
			return nil, err
		}
	}

	return pbUserInfoList, nil

}

func (r *relationRepo) FollowerUserIdList(ctx context.Context, userId int64) ([]int64, error) {
	//1.缓存获取
	var userIds []int64

	ids, err := r.data.rdb.SMembers(ctx, "follower::"+strconv.Itoa(int(userId))).Result()
	if err == redis.Nil {
		//2.数据库获取
		r.data.db.Model(&model.Relation{}).Select("user_id").Where("to_user_id = ?", userId).Find(&userIds)

		//更新缓存
		err := r.data.rdb.SAdd(ctx, "follower::"+strconv.Itoa(int(userId)), userIds).Err()
		if err != nil {
			return nil, err
		}

		return userIds, err

	} else if err != nil {
		return nil, err
	}

	for _, id := range ids {
		userId, _ := strconv.ParseInt(id, 10, 64)
		userIds = append(userIds, userId)
	}

	return userIds, nil
}

// 3.rpc  获取好友最新消息和消息类型
//
//	message  LatestMessage {
//	 string content = 1; // 消息内容
//	 int64 msgType = 2; // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
//	 int64 friend_id = 4; // 好友的id
//	}
func (r *relationRepo) GetNewMessages(ctx context.Context, userId int64, friendIds []int64) ([]*messageV1.LatestMessage, error) {
	//1.获取最新消息
	latestMessages, err := r.data.messagec.GetNewMessages(ctx, &messageV1.GetNewMessagesRequest{UserId: strconv.FormatInt(userId, 10), ToUserId: friendIds})
	if err != nil {
		return nil, err
	}

	return latestMessages.LatestMessageList, nil
}
