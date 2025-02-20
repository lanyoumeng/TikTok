package data

import (
	"context"
	pb "message/api/message/v1"
	"message/internal/biz"
	"message/internal/pkg/model"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type messageRepo struct {
	data *Data
	log  *log.Helper
}

// NewMessageRepo .
func NewMessageRepo(data *Data, logger log.Logger) biz.MessageRepo {
	return &messageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *messageRepo) SaveMessage(ctx context.Context, message model.Message) error {
	start := time.Now()
	// 保存消息
	err := r.data.db.Create(&message).Error
	if err != nil {
		r.log.Error("SaveMessage err:", err)
		return err
	}

	r.log.Infof("SaveMessage success , SaveMessage耗时=%v", time.Since(start))
	return nil
}

func (r *messageRepo) GetMessageRecord(ctx context.Context, userId, toUserId, preMsgTime int64) ([]*model.Message, error) {
	start := time.Now()
	// 获取消息记录
	var messageList []*model.Message

	//preMsgTime字段还没有用到
	//err := r.data.db.Where("from_user_id = ? AND to_user_id = ? AND create_time < ?", userId, toUserId, preMsgTime).Order("create_time desc").Limit(10).Find(&messageList).Error
	err := r.data.db.Where("from_user_id = ? AND to_user_id = ? ", userId, toUserId).Order("create_time desc").Find(&messageList).Error

	if err != nil {
		r.log.Error("GetMessageRecord err:", err)
		return nil, err
	}

	r.log.Infof("GetMessageRecord success , GetMessageRecord耗时=%v", time.Since(start))
	return messageList, nil
}

func (r *messageRepo) GetLatestMessage(ctx context.Context, userId, friendId int64) (*pb.LatestMessage, error) {
	start := time.Now()
	// 获取最新消息 1条
	//看是用户发的，还是好友发的
	var cnt1, cnt2 int64
	var message1, message2 model.Message
	latestMessage := &pb.LatestMessage{}

	//查询好友发送的消息
	err := r.data.db.Where("from_user_id = ? AND to_user_id = ? ", friendId, userId).Order("create_time desc").Count(&cnt1).First(&message1).Error
	if err != nil {
		r.log.Errorf("GetLatestMessage err: %v", err)
		return nil, err
	}

	//查询用户发送的消息
	err = r.data.db.Where("from_user_id = ? AND to_user_id = ? ", userId, friendId).Order("create_time desc").Count(&cnt2).First(&message2).Error
	if err != nil {
		r.log.Errorf("GetLatestMessage err: %v", err)
		return nil, err
	}

	//比较两个消息的时间，返回最新的消息
	if cnt1 > 0 && cnt2 > 0 {
		if message1.CreateTime.After(message2.CreateTime) {
			latestMessage.FriendId = friendId
			latestMessage.Content = message1.Content
			latestMessage.MsgType = 0 //当前请求用户接收的消息
		} else {

			latestMessage.FriendId = friendId
			latestMessage.Content = message2.Content
			latestMessage.MsgType = 1 //当前请求用户发送的消息
		}
	} else if cnt1 > 0 && cnt2 == 0 {

		latestMessage.FriendId = friendId
		latestMessage.Content = message1.Content
		latestMessage.MsgType = 0 //当前请求用户接收的消息

	} else if cnt2 > 0 && cnt1 == 0 {
		latestMessage.FriendId = friendId
		latestMessage.Content = message2.Content
		latestMessage.MsgType = 1 //当前请求用户发送的消息

	} else {
		return nil, nil

	}

	r.log.Infof("GetLatestMessage success , GetLatestMessage耗时=%v", time.Since(start))
	return latestMessage, nil
}
