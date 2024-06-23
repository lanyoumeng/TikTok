package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"

	messageV1 "relation/api/message/v1"
	pb "relation/api/relation/v1"
	"relation/internal/pkg/errno"
)

// RelationRepo is a Greater repo.
type RelationRepo interface {
	IsFollow(ctx context.Context, userId, toUserId int64) (bool, error)
	InsertFollow(ctx context.Context, userId, toUserId int64) error
	DeleteFollow(ctx context.Context, userId, toUserId int64) error

	FollowUserIdList(ctx context.Context, userId int64) ([]int64, error)
	UserInfoList(ctx context.Context, userIdList []int64) ([]*pb.User, error)

	FollowerUserIdList(ctx context.Context, userId int64) ([]int64, error)

	GetNewMessages(ctx context.Context, userId int64, friendIds []int64) ([]*messageV1.LatestMessage, error)
}

// RelationUsecase is a Greeter usecase.
type RelationUsecase struct {
	repo RelationRepo
	log  *log.Helper
}

// NewRelationUsecase new a Greeter usecase.
func NewRelationUsecase(repo RelationRepo, logger log.Logger) *RelationUsecase {
	return &RelationUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (r *RelationUsecase) Follow(ctx context.Context, userId, targetUserId, actionType int64) error {

	//if actionType == 1 关注操作

	if actionType == 1 {
		// 1. 判断是否已经关注
		isFollow, err := r.repo.IsFollow(ctx, userId, targetUserId)
		if err != nil {
			return err
		}
		// 2. 如果已经关注，返回错误
		if isFollow {
			return errno.Errhavefollowed

		}
		// 3. 如果没有关注，插入关注关系
		err = r.repo.InsertFollow(ctx, userId, targetUserId)
		if err != nil {
			return err
		}
	}

	//if actionType == 2 取消关注操作

	if actionType == 2 {
		// 1. 判断是否已经关注
		isFollow, err := r.repo.IsFollow(ctx, userId, targetUserId)
		if err != nil {
			return err
		}

		// 2. 如果没有关注，返回错误
		if !isFollow {
			return errno.Errnotfollowed

		}
		// 3. 如果已经关注，删除关注关系
		err = r.repo.DeleteFollow(ctx, userId, targetUserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RelationUsecase) FollowList(ctx context.Context, userId int64) ([]*pb.User, error) {
	// 1. 获取关注列表
	//获取关注用户的ids 然后rpc获取用户信息
	//再获取每个用户的is_follow字段 即登录用户是否关注了该用户

	followUserIdList, err := r.repo.FollowUserIdList(ctx, userId)
	if err != nil {
		return nil, err
	}

	userInfoList := make([]*pb.User, 0)
	userInfoList, err = r.repo.UserInfoList(ctx, followUserIdList)
	if err != nil {
		return nil, err
	}

	//for _, user := range userInfoList {
	//	//再获取每个用户的is_follow字段
	//	flag, err := r.repo.IsFollow(ctx, user.Id, userId)
	//	if err != nil {
	//		return nil, err
	//	}
	//	user.IsFollow = flag
	//}

	for _, user := range userInfoList {
		user.IsFollow = true
	}

	return userInfoList, nil
}

func (r *RelationUsecase) FollowerList(ctx context.Context, userId int64) ([]*pb.User, error) {
	// 1. 获取粉丝列表
	//获取粉丝用户的ids 然后rpc获取用户信息
	//再获取每个用户的is_follow字段 即登录用户是否关注了该用户

	followerUserIdList, err := r.repo.FollowerUserIdList(ctx, userId)
	if err != nil {
		return nil, err
	}

	userInfoList := make([]*pb.User, 0)
	userInfoList, err = r.repo.UserInfoList(ctx, followerUserIdList)
	if err != nil {
		return nil, err
	}

	for _, user := range userInfoList {
		//再获取每个用户的is_follow字段 登录用户是否关注了该用户
		flag, err := r.repo.IsFollow(ctx, userId, user.Id)
		if err != nil {
			return nil, err
		}
		user.IsFollow = flag
	}

	return userInfoList, nil
}

// 1.获取好友id列表
// 2.获取好友信息
// 3.获取好友最新消息和消息类型
func (r *RelationUsecase) FriendList(ctx context.Context, userId int64) ([]*pb.FriendUser, error) {

	// 1. 获取好友列表
	//好友列表，关注和粉丝的交集

	followUserIdList, err := r.repo.FollowUserIdList(ctx, userId)
	if err != nil {
		return nil, err
	}

	followerUserIdList, err := r.repo.FollowerUserIdList(ctx, userId)
	if err != nil {
		return nil, err
	}

	friendUserIdList := make([]int64, 0)
	for _, followUserId := range followUserIdList {
		for _, followerUserId := range followerUserIdList {
			if followUserId == followerUserId {
				friendUserIdList = append(friendUserIdList, followUserId)
				break
			}
		}
	}

	//2.获取好友信息
	//rpc获取用户信息
	//再rpc获取每个用户的is_follow字段 即登录用户是否关注了该用户

	userInfoList := make([]*pb.User, 0)
	userInfoList, err = r.repo.UserInfoList(ctx, friendUserIdList)
	if err != nil {
		return nil, err
	}

	friendList := make([]*pb.FriendUser, 0)
	for _, user := range userInfoList {
		friendList = append(friendList, &pb.FriendUser{User: user})
	}

	for _, friend := range friendList {
		//再获取每个用户的is_follow字段 登录用户是否关注了该用户
		flag, err := r.repo.IsFollow(ctx, userId, friend.User.Id)
		if err != nil {
			return nil, err
		}
		friend.User.IsFollow = flag
	}

	//3.rpc  获取好友最新消息和消息类型
	//message  LatestMessage {
	//  string content = 1; // 消息内容
	//  int64 msgType = 2; // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
	//  int64 friend_id = 4; // 好友的id
	//}

	newMessages, err := r.repo.GetNewMessages(ctx, userId, friendUserIdList)
	if err != nil {
		return nil, err
	}

	for _, message := range newMessages {
		for _, friend := range friendList {
			if message.FriendId == friend.User.Id {
				friend.Message = message.Content
				friend.MsgType = message.MsgType
			}
			break
		}
	}

	return friendList, nil
}

func (r *RelationUsecase) IsFollow(ctx context.Context, userId, toUserId int64) (bool, error) {
	// 1. 判断是否已经关注
	isFollow, err := r.repo.IsFollow(ctx, userId, toUserId)
	if err != nil {
		return false, err
	}

	return isFollow, nil

}

func (r *RelationUsecase) FollowCnt(ctx context.Context, userId int64) (int64, int64, error) {
	// 1. 获取关注数和粉丝数
	followUserIdList, err := r.repo.FollowUserIdList(ctx, userId)
	if err != nil {
		return 0, 0, err
	}

	followerUserIdList, err := r.repo.FollowerUserIdList(ctx, userId)
	if err != nil {
		return 0, 0, err
	}

	return int64(len(followUserIdList)), int64(len(followerUserIdList)), nil
}
