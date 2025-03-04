package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"relation/internal/biz"
	"relation/internal/conf"
	"relation/pkg/token"
	"time"

	pb "relation/api/relation/v1"
)

type RelationService struct {
	pb.UnimplementedRelationServiceServer
	ru     *biz.RelationUsecase
	JwtKey string
	log    *log.Helper
}

func NewRelationService(uc *biz.RelationUsecase, auth *conf.Auth, logger log.Logger) *RelationService {
	return &RelationService{ru: uc, JwtKey: auth.JwtKey, log: log.NewHelper(logger)}
}

func (s *RelationService) Relation(ctx context.Context, req *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	start := time.Now()
	user, err := token.ParseToken(req.Token, s.JwtKey)
	if err != nil {
		s.log.Errorf("token.ParseToken error: %v", err)
		return &pb.DouyinRelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "token解析失败",
		}, err
	}

	err = s.ru.Follow(ctx, user.UserId, req.ToUserId, int64(req.ActionType))
	if err != nil {
		s.log.Errorf("Follow error: %v", err)
		return &pb.DouyinRelationActionResponse{
			StatusCode: 1,
			StatusMsg:  "操作失败",
		}, err

	}

	s.log.Infof("Relation end , Relationh耗时:%v", time.Since(start))
	return &pb.DouyinRelationActionResponse{
		StatusCode: 0,
		StatusMsg:  "操作成功",
	}, nil
}
func (s *RelationService) RelationFollowList(ctx context.Context, req *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {
	start := time.Now()
	//关注列表
	//获取user
	userInfoList := make([]*pb.User, 5)

	userInfoList, err := s.ru.FollowList(ctx, req.UserId)
	if err != nil {
		s.log.Errorf("FollowList error: %v", err)
		return nil, err
	}

	s.log.Infof("RelationFollowList end ，RelationFollowList耗时:%v", time.Since(start))
	return &pb.DouyinRelationFollowListResponse{
		StatusCode: 0,
		StatusMsg:  "获取关注列表成功",
		UserList:   userInfoList,
	}, nil
}
func (s *RelationService) RelationFollowerList(ctx context.Context, req *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {
	start := time.Now()
	//粉丝列表
	//获取user
	userInfoList := make([]*pb.User, 5)

	userInfoList, err := s.ru.FollowerList(ctx, req.UserId)
	if err != nil {
		s.log.Errorf("FollowerList error: %v", err)
		return nil, err
	}

	s.log.Infof("RelationFollowerList end ，RelationFollowerList耗时:%v", time.Since(start))
	return &pb.DouyinRelationFollowerListResponse{
		StatusCode: 0,
		StatusMsg:  "获取粉丝列表成功",
		UserList:   userInfoList,
	}, nil
}
func (s *RelationService) FriendList(ctx context.Context, req *pb.DouyinRelationFriendListRequest) (*pb.DouyinRelationFriendListResponse, error) {
	start := time.Now()
	//好友列表，关注和粉丝的交集
	//获取user
	friendList := make([]*pb.FriendUser, 5)

	//message FriendUser {
	//  User user = 1; // 嵌套User消息
	//  string message = 2; // 和该好友的最新聊天消息
	//  int64 msgType = 3; // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	//}
	//1.获取好友id列表
	//2.获取好友信息
	//3.获取好友最新消息和消息类型

	friendList, err := s.ru.FriendList(ctx, req.UserId)
	if err != nil {
		s.log.Errorf("FriendList error: %v", err)
		return nil, err
	}

	s.log.Infof("FriendList end ，FriendList耗时:%v", time.Since(start))
	return &pb.DouyinRelationFriendListResponse{
		StatusCode: 0,
		StatusMsg:  "获取朋友列表成功",
		UserList:   friendList,
	}, nil
}

// 获取关注数和粉丝数
func (s *RelationService) FollowCnt(ctx context.Context, req *pb.FollowCntRequest) (*pb.FollowCntResponse, error) {
	start := time.Now()
	//s.log.Debugf("FollowCnt request: %v", req)
	//获取user
	followCnt, followerCnt, err := s.ru.FollowCnt(ctx, req.UserId)
	if err != nil {
		s.log.Errorf("FollowCnt error: %v", err)
		return nil, err

	}

	s.log.Infof("FollowCnt end ，FollowCnt耗时:%v", time.Since(start))
	return &pb.FollowCntResponse{
		FollowCnt:   followCnt,
		FollowerCnt: followerCnt,
	}, nil
}

// //根据userId,authorId查询用户是否关注作者
// rpc IsFollow (IsFollow_request) returns (IsFollow_response);
func (s *RelationService) IsFollow(ctx context.Context, req *pb.IsFollowRequest) (*pb.IsFollowResponse, error) {

	start := time.Now()
	flag, err := s.ru.IsFollow(ctx, req.UserId, req.AuthorId)
	if err != nil {
		s.log.Errorf("IsFollow error: %v", err)
		return nil, err
	}

	s.log.Infof("IsFollow end ，IsFollow耗时:%v", time.Since(start))
	return &pb.IsFollowResponse{
		IsFollow: flag,
	}, nil
}
