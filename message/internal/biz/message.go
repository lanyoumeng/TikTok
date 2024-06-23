package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	pb "message/api/message/v1"
	"message/internal/pkg/model"
	"time"
)

// MessageRepo is a Greater repo.
type MessageRepo interface {
	SaveMessage(ctx context.Context, message model.Message) error
	GetMessageRecord(ctx context.Context, userId, toUserId, preMsgTime int64) ([]*model.Message, error)

	GetLatestMessage(ctx context.Context, userId, friendId int64) (*pb.LatestMessage, error)
}

// MessageUsecase is a Greeter usecase.
type MessageUsecase struct {
	repo MessageRepo
	log  *log.Helper
}

// NewMessageUsecase new a Greeter usecase.
func NewMessageUsecase(repo MessageRepo, logger log.Logger) *MessageUsecase {
	return &MessageUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (m *MessageUsecase) SendMessage(ctx context.Context, message model.Message) error {
	// 发送消息

	err := m.repo.SaveMessage(ctx, message)
	if err != nil {
		return err
	}
	return nil

}

func (m *MessageUsecase) GetMessageRecord(ctx context.Context, userId, toUserId, preMsgTime int64) ([]*pb.Message, error) {
	// 获取消息记录
	// 1. 查询mysql消息记录
	messageList, err := m.repo.GetMessageRecord(ctx, userId, toUserId, preMsgTime)
	if err != nil {
		return nil, err
	}
	// 2. 转换为pb.Message
	var pbMessageList []*pb.Message
	for _, message := range messageList {
		// 将 CreatedAt 转换为字符串
		createdAtStr := message.CreatedAt.Format(time.RFC3339)
		pbMessage := &pb.Message{
			Id:         message.Id,
			FromUserId: message.FromUserId,
			ToUserId:   message.ToUserId,
			Content:    message.Content,
			CreateTime: createdAtStr,
		}
		pbMessageList = append(pbMessageList, pbMessage)
	}
	return pbMessageList, nil

}

// MessageRecordRequest .
func (m *MessageUsecase) GetLatestMessage(ctx context.Context, userId int64, friendIds []int64) ([]*pb.LatestMessage, error) {
	// 获取最新消息
	// 1. 查询mysql消息记录
	var latestMessageList []*pb.LatestMessage
	for _, friendId := range friendIds {
		// 获取最新消息 1条
		//看是用户发的，还是好友发的
		latestMessage, err := m.repo.GetLatestMessage(ctx, userId, friendId)
		if err != nil {
			return nil, err
		}

		latestMessageList = append(latestMessageList, latestMessage)

	}
	return latestMessageList, nil
}
