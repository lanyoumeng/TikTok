package service

import (
	"context"
	"fmt"
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
	start := time.Now()
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

	////将feedVideo转换为feed
	//feed := make([]pb.Video, len(feedVideo))
	//
	//for k, v := range feedVideo {
	//	feed[k] = *v
	//}

	v.log.Infof("service.Feed success , Feed耗时=%v", time.Since(start))

	return &pb.DouyinFeedResponse{
		StatusCode: 0,
		StatusMsg:  "返回视频流成功",
		VideoList:  feedVideo,
		NextTime:   earliestTime,
	}, nil
}
func (v *VideoService) Publish(ctx context.Context, req *pb.DouyinPublishActionRequest) (*pb.DouyinPublishActionResponse, error) {
	start := time.Now()
	fmt.Printf("title:%v ,data:%v, token:%v", req.Title, req.Data, req.Token)
	//获取user
	user, err := token.ParseToken(req.Token, v.JwtKey)
	if err != nil {
		log.Debug("service.Publish/ParseToken", err)
		return nil, err
	}
	err = v.vc.Publish(ctx, req.Title, &req.Data, user.UserId)
	if err != nil {
		log.Debug("service.Publish/Publish", err)
		return nil, err

	}

	v.log.Infof("service.Publish success , Publish耗时=%v", time.Since(start))
	return &pb.DouyinPublishActionResponse{
		StatusCode: 0,
		StatusMsg:  "发布视频成功",
	}, nil
}
func (v *VideoService) PublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (*pb.DouyinPublishListResponse, error) {
	start := time.Now()

	////获取user
	//user, err := token.ParseToken(req.Token, v.JwtKey)
	//if err != nil {
	//	return nil, err
	//}
	userId := req.UserId
	publishList, err := v.vc.PublishList(ctx, userId)
	if err != nil {
		return nil, err
	}

	v.log.Infof("service.PublishList success , PublishList耗时=%v", time.Since(start))
	return &pb.DouyinPublishListResponse{
		StatusCode: 0,
		StatusMsg:  "获取发布视频列表成功",
		VideoList:  publishList,
	}, nil
}

func (v *VideoService) WorkCnt(ctx context.Context, req *pb.WorkCntRequest) (*pb.WorkCntResponse, error) {
	start := time.Now()

	workCount, err := v.vc.WorkCnt(ctx, req.UserId)
	if err != nil {
		log.Debug("workCount", workCount, "err", err)
		return nil, err
	}

	v.log.Infof("service.WorkCnt success , WorkCnt耗时=%v", time.Since(start))
	return &pb.WorkCntResponse{
		WorkCount: workCount,
	}, nil
}

func (v *VideoService) FavoriteListByVId(ctx context.Context, req *pb.FavoriteListReq) (*pb.FavoriteListResp, error) {
	start := time.Now()

	videoList, err := v.vc.FavoriteListByVId(ctx, req.VideoIdList)
	if err != nil {
		return nil, err
	}

	v.log.Infof("service.FavoriteListByVId success , FavoriteListByVId耗时=%v", time.Since(start))
	return &pb.FavoriteListResp{
		VideoList: videoList,
	}, nil
}

// //通过作者id 获取作者的发布视频id列表
// rpc PublishVidsByAId(PublishVidsByAId_req) returns (PublishVidsByAIdResp);
func (v *VideoService) PublishVidsByAId(ctx context.Context, req *pb.PublishVidsByAIdReq) (*pb.PublishVidsByAIdResp, error) {
	start := time.Now()
	videoIds, err := v.vc.PublishVidsByAId(ctx, req.AuthorId)
	if err != nil {
		return nil, err
	}

	v.log.Infof("service.PublishVidsByAId success , PublishVidsByAId耗时=%v", time.Since(start))
	return &pb.PublishVidsByAIdResp{
		VideoIdList: videoIds,
	}, nil

}

// 通过视频id获取作者id
func (v *VideoService) GetAIdByVId(ctx context.Context, req *pb.GetAIdByVIdReq) (*pb.GetAIdByVIdResp, error) {
	start := time.Now()
	authorId, err := v.vc.GetAIdByVId(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}

	v.log.Infof("service.GetAIdByVId success , GetAIdByVId耗时=%v", time.Since(start))
	return &pb.GetAIdByVIdResp{
		AuthorId: authorId,
	}, nil
}
