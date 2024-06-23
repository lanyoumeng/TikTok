package service

import (
	pb "comment/api/comment/v1"
	"comment/internal/biz"
	"comment/internal/conf"
	"comment/internal/pkg/model"
	"comment/pkg/token"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type CommentService struct {
	pb.UnimplementedCommentServiceServer
	uc     *biz.CommentUsecase
	JwtKey string
	log    *log.Helper
}

func NewCommentService(uc *biz.CommentUsecase, auth *conf.Auth, logger log.Logger) *CommentService {

	return &CommentService{uc: uc, JwtKey: auth.JwtKey, log: log.NewHelper(logger)}
}

// rpc Comment (douyin_comment_send_request) returns(douyin_comment_send_response);
func (c *CommentService) Comment(ctx context.Context, req *pb.DouyinCommentSendRequest) (*pb.DouyinCommentSendResponse, error) {
	user, err := token.ParseToken(req.Token, c.JwtKey)
	if err != nil {
		return nil, err
	}
	// 1-发布评论，2-删除评论
	// 发布评论
	if req.SendType == 1 {
		comment := &model.Comment{}
		comment.Content = req.CommentText
		comment.UserId = user.UserId
		comment.VideoId = req.VideoId

		comment, err := c.uc.SendComment(ctx, comment)
		if err != nil {
			return nil, err
		}

		commentresp := &pb.Comment{}
		// 评论用户信息
		userinfo := &pb.User{}
		// getuserinfoByUIdVIdAId
		userinfo, err = c.uc.GetUserinfoByUIdVIdAId(ctx, user.UserId, req.VideoId)
		if err != nil {
			return nil, err
		}
		commentresp.Id = comment.Id
		commentresp.Content = comment.Content
		commentresp.User = userinfo
		commentresp.CreateDate = comment.CreateDate

		return &pb.DouyinCommentSendResponse{
			StatusCode: 0,
			StatusMsg:  "push comment success",
			Comment:    commentresp,
		}, nil

	}

	//删除评论
	if req.SendType == 1 {

		err := c.uc.DelComment(ctx, req.CommentId)
		if err != nil {
			return nil, err
		}

	}

	return &pb.DouyinCommentSendResponse{
		StatusCode: 0,
		StatusMsg:  "Del comment success",
	}, nil
}

//// 获取评论列表 注意每个评论用户信息中的is_follow字段需要从favorite服务单独获取
//rpc CommentList(douyin_comment_list_request) returns (douyin_comment_list_response);

func (c *CommentService) CommentList(ctx context.Context, req *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {

	// 获取评论列表
	commentList, err := c.uc.CommentList(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}

	return &pb.DouyinCommentListResponse{
		StatusCode:  0,
		StatusMsg:   "get comment list success",
		CommentList: commentList,
	}, nil
}

// rpc GetCommentCntByVId(GetCommentCntByVIdReq) returns (GetCommentCntByVIdResp);
func (c *CommentService) GetCommentCntByVId(ctx context.Context, req *pb.GetCommentCntByVIdReq) (*pb.GetCommentCntByVIdResp, error) {
	//GetCommentCntByVId
	cnt, err := c.uc.GetCommentCntByVId(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}

	return &pb.GetCommentCntByVIdResp{
		CommentCount: cnt,
	}, nil
}
