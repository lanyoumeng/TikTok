package service

import (
	"context"
	"favorite/pkg/tool"
	"time"

	pb "favorite/api/favorite/v1"
	"favorite/internal/biz"
	"favorite/internal/conf"
	"favorite/pkg/token"

	"github.com/go-kratos/kratos/v2/log"
)

type FavoriteService struct {
	pb.UnimplementedFavoriteServiceServer
	fc     *biz.FavoriteUsecase
	JwtKey string
	log    *log.Helper
}

func NewFavoriteService(fc *biz.FavoriteUsecase, auth *conf.Auth, logger log.Logger) *FavoriteService {
	return &FavoriteService{fc: fc, JwtKey: auth.JwtKey, log: log.NewHelper(logger)}
}

// 点赞或者取消点赞
func (f *FavoriteService) Favorite(ctx context.Context, req *pb.DouyinFavoriteActionRequest) (*pb.DouyinFavoriteActionResponse, error) {
	start := time.Now()
	//获取user
	user, err := token.ParseToken(req.Token, f.JwtKey)
	if err != nil {
		f.log.Errorf("token.ParseToken error: %v", err)
		return nil, err
	}
	//获取作者id
	authorId, err := f.fc.GetAuthorIdByVideoId(ctx, req.VideoId)
	if err != nil {
		f.log.Errorf("GetAuthorIdByVideoId error: %v", err)
		return nil, err
	}
	if err := f.fc.Favorite(ctx, authorId, req.VideoId, user.UserId, int64(req.ActionType)); err != nil {
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "FavoriteAction fale",
		}, nil

	}

	f.log.Infof("service.Favorite success , Favorite耗时=%v", time.Since(start))

	return &pb.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "FavoriteAction success",
	}, nil
}

func (f *FavoriteService) FavoriteList(ctx context.Context, req *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {
	start := time.Now()
	videoList, err := f.fc.FavoriteList(ctx, req.UserId)
	if err != nil {
		f.log.Errorf("FavoriteList error: %v", err)
		return nil, err
	}
	videos := make([]*pb.Video, len(videoList))

	videos = tool.ConvertVideoList(videoList)

	f.log.Infof("service.FavoriteList success , FavoriteList耗时=%v", time.Since(start))
	return &pb.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "FavoriteList success",
		VideoList:  videos,
	}, nil
}
func (f *FavoriteService) GetFavoriteCntByVId(ctx context.Context, req *pb.GetFavoriteCntByVIdRequest) (*pb.GetFavoriteCntByVIdResponse, error) {
	start := time.Now()
	cnt, err := f.fc.GetFavoriteCntByVId(ctx, req.Id)
	if err != nil {
		f.log.Errorf("GetFavoriteCntByVId error: %v", err)
		return nil, err
	}

	f.log.Infof("service.GetFavoriteCntByVId success , GetFavoriteCntByVId耗时=%v", time.Since(start))
	return &pb.GetFavoriteCntByVIdResponse{
		FavoriteCount: cnt,
	}, nil
}
func (f *FavoriteService) GetIsFavorite(ctx context.Context, req *pb.GetIsFavoriteRequest) (*pb.GetIsFavoriteResponse, error) {
	start := time.Now()
	flag, err := f.fc.GetIsFavorite(ctx, req.VideoId, req.UserId)
	if err != nil {
		f.log.Errorf("GetIsFavorite error: %v", err)
		return nil, err
	}

	f.log.Infof("service.GetIsFavorite success , GetIsFavorite耗时=%v", time.Since(start))
	return &pb.GetIsFavoriteResponse{
		Favorite: flag,
	}, nil
}

// 获取 用户的 获赞数TotalFavorited 和 点赞数量FavoriteCount
func (f *FavoriteService) GetFavoriteCntByUId(ctx context.Context, req *pb.GetFavoriteCntByUIdRequest) (*pb.GetFavoriteCntByUIdResponse, error) {
	start := time.Now()
	TotalFavorited, FavoriteCount, err := f.fc.GetFavoriteCntByUId(ctx, req.UserId)
	if err != nil {
		f.log.Errorf("GetFavoriteCntByUId error: %v", err)
		return nil, err
	}

	f.log.Infof("service.GetFavoriteCntByUId success , GetFavoriteCntByUId耗时=%v", time.Since(start))
	return &pb.GetFavoriteCntByUIdResponse{
		TotalFavorited: TotalFavorited,
		FavoriteCount:  FavoriteCount,
	}, nil
}
