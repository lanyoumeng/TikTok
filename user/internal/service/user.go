package service

import (
	"context"
	"strconv"
	pb "user/api/user/v1"
	"user/internal/biz"
	"user/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	uc  *biz.UserUsecase
	log *log.Helper
}

func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

func (s *UserService) Register(ctx context.Context, req *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	// var user *biz.User
	user := &model.User{} // 初始化 user
	user.Name = req.Username
	user.Password = req.Password

	//s.log.Debug("user:", user)
	userId, token, err := s.uc.Create(ctx, user)
	if err != nil {
		s.log.Error("service.Register/Create-err:", err)
		return &pb.UserRegisterResponse{
			StatusCode: -1,
			StatusMsg:  "创建用户失败/用户已存在",
		}, err
	}

	return &pb.UserRegisterResponse{
		StatusCode: 0,
		StatusMsg:  "创建用户成功",
		UserId:     userId,
		Token:      token,
	}, nil
}
func (s *UserService) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {

	// var user *biz.User
	user := &model.User{} // 初始化 user
	user.Name = req.Username
	user.Password = req.Password

	userId, token, err := s.uc.Login(ctx, user)
	if err != nil {
		s.log.Error("service.Login/Login-err:", err)

		return &pb.UserLoginResponse{
			StatusCode: -1,
			StatusMsg:  "登录失败",
		}, err

	}
	return &pb.UserLoginResponse{
		StatusCode: 0,
		StatusMsg:  "登录成功",
		UserId:     userId,
		Token:      token}, nil
}

func (s *UserService) UserInfo(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {

	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return &pb.UserResponse{
			StatusCode: -1,
			StatusMsg:  "ParseInt函数失败",
		}, err

	}
	user, err := s.uc.UserInfo(ctx, userId)
	if err != nil {
		s.log.Error("err:", err)
		return &pb.UserResponse{
			StatusCode: -1,
			StatusMsg:  "获取用户信息失败",
		}, err
	}

	return &pb.UserResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户信息成功",
		User:       user,
	}, nil
}

// //WorkCount
// rpc UpdateWorkCnt(UpdateWorkCntRequest) returns(UpdateWorkCntResponse);
func (s *UserService) UpdateWorkCnt(ctx context.Context, req *pb.UpdateWorkCntRequest) (*pb.UpdateWorkCntResponse, error) {
	userId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	count := &model.UserCount{}
	count, err = s.uc.RGetCountById(ctx, userId)
	if err != nil {
		s.log.Errorf("err:%v  userId:%v", err, userId)
		return nil, err
	}
	count.WorkCount = req.WorkCount

	err = s.uc.RSaveCount(ctx, count)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err

	}

	return &pb.UpdateWorkCntResponse{
		StatusCode: 0,
		StatusMsg:  "作品数量更新",
	}, nil

}

func (s *UserService) UpdateFavoriteCnt(ctx context.Context, req *pb.UpdateFavoriteCntRequest) (*pb.UpdateFavoriteCntResponse, error) {

	//更新点赞用户的计数信息 FavoriteUserId
	FavoriteUserId, err := strconv.ParseInt(req.FavoriteUserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}
	count := &model.UserCount{}
	count, err = s.uc.RGetCountById(ctx, FavoriteUserId)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}
	count.FavoriteCount = req.FavoriteCount

	err = s.uc.RSaveCount(ctx, count)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	//更新被点赞用户的计数信息 FavoritedUserId
	FavoritedUserId, err := strconv.ParseInt(req.FavoritedUserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err

	}
	count, err = s.uc.RGetCountById(ctx, FavoritedUserId)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err

	}
	count.TotalFavorited = req.TotalFavorited

	err = s.uc.RSaveCount(ctx, count)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	return &pb.UpdateFavoriteCntResponse{
		StatusCode: 0,
		StatusMsg:  "点赞计数更新",
	}, nil
}

func (s *UserService) UpdateFollowCnt(ctx context.Context, req *pb.UpdateFollowCntRequest) (*pb.UpdateFollowCntResponse, error) {
	//更新执行关注的 用户的计数信息 FollowUserId
	userId, err := strconv.ParseInt(req.FollowUserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	count := &model.UserCount{}
	count, err = s.uc.RGetCountById(ctx, userId)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}
	count.FollowCount = req.FollowCount

	err = s.uc.RSaveCount(ctx, count)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	//更新被关注用户的计数信息 FollowedUserId
	followedUserId, err := strconv.ParseInt(req.FollowedUserId, 10, 64)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	count, err = s.uc.RGetCountById(ctx, followedUserId)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	count.FollowerCount = req.FollowerCount

	err = s.uc.RSaveCount(ctx, count)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	return &pb.UpdateFollowCntResponse{
		StatusCode: 0,
		StatusMsg:  "关注计数更新",
	}, nil
}

// 用户信息列表 uIds
// rpc UserInfoList (UserInfoListrRequest) returns (UserInfoListResponse);
func (s *UserService) UserInfoList(ctx context.Context, req *pb.UserInfoListrRequest) (*pb.UserInfoListResponse, error) {

	//获取用户信息列表
	userInfoList := make([]*pb.User, 5)
	userInfoList, err := s.uc.UserInfoList(ctx, req.UserId)
	if err != nil {
		s.log.Error("err:", err)
		return nil, err
	}

	return &pb.UserInfoListResponse{
		StatusCode: 0,
		StatusMsg:  "获取用户信息列表成功",
		Users:      userInfoList,
	}, nil
}
