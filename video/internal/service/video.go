package service

import (
	"context"
	"time"
	"video/pkg/token"

	pb "video/api/video/v1"
	"video/internal/biz"
	"video/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

type VideoService struct {
	pb.UnimplementedVideoServiceServer

	vc     *biz.VideoUsecase
	JwtKey string
	log    *log.Helper
}

func NewVideoService(vc *biz.VideoUsecase, auth *conf.Auth, logger log.Logger) *VideoService {
	return &VideoService{vc: vc, JwtKey: auth.JwtKey, log: log.NewHelper(logger)}
}

func (v *VideoService) Feed(ctx context.Context, req *pb.DouyinFeedRequest) (*pb.DouyinFeedResponse, error) {

	//lastTime
	latestTime := req.LatestTime
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeParse := time.Unix(latestTime, 0).In(loc)

	if timeParse.After(time.Now()) || timeParse.IsZero() {
		timeParse = time.Now()
	}
	//获取user
	//kratos的方法  Claims, err := jwt.FromContext(ctx)
	user, err := token.ParseToken(req.Token, v.JwtKey)
	if err != nil {
		return nil, err
	}
	feedVideo, earliestTime, err := v.vc.Feed(ctx, timeParse, user.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "返回视频流成功",
		VideoList:  feedVideo,
		NextTime:   earliestTime,
	}, nil
}
func (v *VideoService) Publish(ctx context.Context, req *pb.DouyinPublishActionRequest) (*pb.DouyinPublishActionResponse, error) {
	//获取user
	user, err := token.ParseToken(req.Token, v.JwtKey)
	if err != nil {
		return nil, err
	}
	err = v.vc.Publish(ctx, req.Title, &req.Data, user.UserId)

	return &pb.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "发布视频成功",
	}, nil
}
func (v *VideoService) PublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (*pb.DouyinPublishListResponse, error) {

	//获取user
	user, err := token.ParseToken(req.Token, v.JwtKey)
	if err != nil {
		return nil, err
	}
	publishList, err := v.vc.PublishList(ctx, user.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取发布视频列表成功",
		VideoList:  publishList,
	}, nil
}

func (v *VideoService) WorkCnt(ctx context.Context, req *pb.WorkCntRequest) (*pb.WorkCntResponse, error) {

	workCount, err := v.vc.WorkCnt(ctx, req.UserId)
	if err != nil {
		log.Debug("workCount", workCount, "err", err)
		return nil, err
	}
	return &pb.WorkCntResponse{
		WorkCount: workCount,
	}, nil
}

func (v *VideoService) FavoriteListByVId(ctx context.Context, req *pb.FavoriteListReq) (*pb.FavoriteListResp, error) {

	videoList, err := v.vc.FavoriteListByVId(ctx, req.VideoIdList)
	if err != nil {
		return nil, err
	}
	return &pb.FavoriteListResp{
		VideoList: videoList,
	}, nil
}

// //通过作者id 获取作者的发布视频id列表
// rpc PublishVidsByAId(PublishVidsByAId_req) returns (PublishVidsByAIdResp);
func (v *VideoService) PublishVidsByAId(ctx context.Context, req *pb.PublishVidsByAIdReq) (*pb.PublishVidsByAIdResp, error) {
	videoIds, err := v.vc.PublishVidsByAId(ctx, req.AuthorId)
	if err != nil {
		return nil, err
	}
	return &pb.PublishVidsByAIdResp{
		VideoIdList: videoIds,
	}, nil

}

// 通过视频id获取作者id
func (v *VideoService) GetAIdByVId(ctx context.Context, req *pb.GetAIdByVIdReq) (*pb.GetAIdByVIdResp, error) {
	authorId, err := v.vc.GetAIdByVId(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	return &pb.GetAIdByVIdResp{
		AuthorId: authorId,
	}, nil
}
