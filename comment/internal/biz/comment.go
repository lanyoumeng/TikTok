package biz

import (
	pb "comment/api/comment/v1"
	"comment/internal/pkg/model"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

// CommentRepo is a Greater repo.
//
//go:generate mockgen -destination=../mocks/mrepo/comment.go -package=mrepo . CommentRepo
type CommentRepo interface {
	SaveComment(context.Context, *model.Comment) (*model.Comment, error)
	DelComment(context.Context, int64) error

	//根据userId查询用户信息
	GetUserinfoByUId(ctx context.Context, userId int64) (*pb.User, error)
	//根据videoId查询视频Id
	GetAuthorIdByVId(ctx context.Context, videoId int64) (int64, error)
	//根据userId,authorId查询用户是否关注作者
	GetFollowByUIdAId(ctx context.Context, userId, authorId int64) (bool, error)

	//评论列表
	CommentList(ctx context.Context, videoId int64) ([]*model.Comment, error)

	//获取评论数
	GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error)
}

// CommentUsecase is a Greeter usecase.
type CommentUsecase struct {
	repo CommentRepo
	log  *log.Helper
}

// NewCommentUsecase new a Greeter usecase.
func NewCommentUsecase(repo CommentRepo, logger log.Logger) *CommentUsecase {
	return &CommentUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *CommentUsecase) SendComment(ctx context.Context, c *model.Comment) (*model.Comment, error) {
	return uc.repo.SaveComment(ctx, c)
}

func (uc *CommentUsecase) DelComment(ctx context.Context, commentId int64) error {
	return uc.repo.DelComment(ctx, commentId)
}

// 评论用户信息
func (uc *CommentUsecase) GetUserinfoByUIdVIdAId(ctx context.Context, userId, videoId int64) (*pb.User, error) {

	//1.根据userId查询用户信息
	userInfo, err := uc.repo.GetUserinfoByUId(ctx, userId)
	if err != nil {
		uc.log.Errorf("Error getting user info: %v", err)
		return nil, err

	}

	//2.根据videoId查询作者Id
	vId, err := uc.repo.GetAuthorIdByVId(ctx, videoId)
	if err != nil {
		uc.log.Errorf("Error getting author info: %v", err)
		return nil, err

	}

	//3.根据userId,authorId查询用户是否关注作者
	isFollow, err := uc.repo.GetFollowByUIdAId(ctx, userId, vId)
	if err != nil {
		uc.log.Errorf("Error getting follow info: %v", err)
		return nil, err
	}
	userInfo.IsFollow = isFollow

	return userInfo, nil
}

// CommentList 方法
func (uc *CommentUsecase) CommentList(ctx context.Context, videoId int64) ([]*pb.Comment, error) {
	// 获取mysql 评论列表
	comments, err := uc.repo.CommentList(ctx, videoId)
	if err != nil {
		uc.log.Errorf("Error getting comment list: %v", err)
		return nil, err
	}

	//uc.log.Infof("comments: %v", comments)

	var wg sync.WaitGroup
	commentList := make([]*pb.Comment, len(comments))
	errChan := make(chan error, len(comments)) // 错误通道，用于捕获并发中的错误
	defer close(errChan)

	for i, comment := range comments {
		wg.Add(1)
		i, comment := i, comment // 捕获变量
		go func() {
			defer wg.Done()

			// 获取用户信息
			userinfo, err := uc.GetUserinfoByUIdVIdAId(ctx, comment.UserId, videoId)
			if err != nil {
				uc.log.Errorf("Error getting user info: %v", err)
				errChan <- err
				return
			}

			// 构造评论返回
			commentresp := &pb.Comment{
				Id:         comment.Id,
				Content:    comment.Content,
				User:       userinfo,
				CreateDate: comment.CreateDate,
			}

			// 无需加锁，直接赋值
			commentList[i] = commentresp
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 检查是否有错误发生
	select {
	case err := <-errChan:
		// 如果有错误，则返回第一个捕获的错误
		return nil, err
	default:
		// 否则返回评论列表
		return commentList, nil
	}
}

// GetCommentCntByVId
func (uc *CommentUsecase) GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error) {
	return uc.repo.GetCommentCntByVId(ctx, videoId)
}
