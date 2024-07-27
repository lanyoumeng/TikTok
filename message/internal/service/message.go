package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"message/internal/biz"
	"message/internal/conf"
	"message/internal/pkg/model"
	"message/pkg/token"
	"strconv"

	pb "message/api/message/v1"
)

type MessageService struct {
	pb.UnimplementedMessageServiceServer
	uc     *biz.MessageUsecase
	JwtKey string
	log    *log.Helper
}

func NewMessageService(vc *biz.MessageUsecase, auth *conf.Auth, logger log.Logger) *MessageService {
	return &MessageService{uc: vc, JwtKey: auth.JwtKey, log: log.NewHelper(logger)}
}

func (s *MessageService) MessageRecord(ctx context.Context, req *pb.DouyinMessageRecordRequest) (*pb.DouyinMessageRecordResponse, error) {

	//message douyin_message_record_request {
	// string token = 1; // 用户鉴权token
	//  int64 to_user_id = 2; // 对方用户id
	//   int64 pre_msg_time=3;//上次最新消息的时间（新增字段-apk更新中）
	//}

	//获取user
	user, err := token.ParseToken(req.Token, s.JwtKey)
	if err != nil {
		s.log.Errorf("token.ParseToken error: %v", err)
		return nil, err

	}
	//获取消息记录
	messageRecord, err := s.uc.GetMessageRecord(ctx, user.UserId, req.ToUserId, req.PreMsgTime)
	if err != nil {
		s.log.Errorf("GetMessageRecord error: %v", err)
		return nil, err

	}
	return &pb.DouyinMessageRecordResponse{
		StatusCode:  0,
		StatusMsg:   "获取消息记录成功",
		MessageList: messageRecord,
	}, nil
}
func (s *MessageService) MessageSend(ctx context.Context, req *pb.DouyinMessageSendRequest) (*pb.DouyinMessageSendResponse, error) {
	user, err := token.ParseToken(req.Token, s.JwtKey)
	if err != nil {
		s.log.Errorf("token.ParseToken error: %v", err)
		return nil, err
	}
	message := model.Message{
		FromUserId: user.UserId,
		ToUserId:   req.ToUserId,
		Content:    req.Content,
	}

	err = s.uc.SendMessage(ctx, message)

	if err != nil {
		s.log.Errorf("SendMessage error: %v", err)
		return nil, err

	}

	return &pb.DouyinMessageSendResponse{
		StatusCode: 0,
		StatusMsg:  "发送消息成功",
	}, nil
}

//  //  string message = 2; // 和该好友的最新聊天消息
//  //  int64 msgType = 3; // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
//  //获得用户和每个好友的最新一条消息
//  rpc GetNewMessages (GetNewMessages_request) returns (GetNewMessages_response);

func (s *MessageService) GetNewMessages(ctx context.Context, req *pb.GetNewMessagesRequest) (*pb.GetNewMessagesResponse, error) {
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		s.log.Errorf("strconv.ParseInt error: %v", err)
		return nil, err

	}
	latestMessageList, err := s.uc.GetLatestMessage(ctx, userId, req.ToUserId)
	if err != nil {
		s.log.Errorf("GetLatestMessage error: %v", err)
		return nil, err
	}

	return &pb.GetNewMessagesResponse{
		LatestMessageList: latestMessageList,
	}, nil

}
