package service

import (
	"context"
	"favorite/pkg/tool"

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

	//获取user
	user, err := token.ParseToken(req.Token, f.JwtKey)
	if err != nil {
		return nil, err
	}
	//获取作者id
	authorId, err := f.fc.GetAuthorIdByVideoId(ctx, req.VideoId)
	if err != nil {
		return nil, err
	}
	if err := f.fc.Favorite(ctx, authorId, req.VideoId, user.UserId, int64(req.ActionType)); err != nil {
		return &pb.DouyinFavoriteActionResponse{
			StatusCode: -1,
			StatusMsg:  "FavoriteAction fale",
		}, nil

	}

	return &pb.DouyinFavoriteActionResponse{
		StatusCode: 0,
		StatusMsg:  "FavoriteAction success",
	}, nil
}

func (f *FavoriteService) FavoriteList(ctx context.Context, req *pb.DouyinFavoriteListRequest) (*pb.DouyinFavoriteListResponse, error) {

	videoList, err := f.fc.FavoriteList(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	videos := make([]*pb.Video, len(videoList))

	videos = tool.ConvertVideoList(videoList)

	return &pb.DouyinFavoriteListResponse{
		StatusCode: 0,
		StatusMsg:  "FavoriteList success",
		VideoList:  videos,
	}, nil
}
func (f *FavoriteService) GetFavoriteCntByVId(ctx context.Context, req *pb.GetFavoriteCntByVIdRequest) (*pb.GetFavoriteCntByVIdResponse, error) {

	cnt, err := f.fc.GetFavoriteCntByVId(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetFavoriteCntByVIdResponse{
		FavoriteCount: cnt,
	}, nil
}
func (f *FavoriteService) GetIsFavorite(ctx context.Context, req *pb.GetIsFavoriteRequest) (*pb.GetIsFavoriteResponse, error) {

	flag, err := f.fc.GetIsFavorite(ctx, req.VideoId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetIsFavoriteResponse{
		Favorite: flag,
	}, nil
}

// 获取 用户的 获赞数TotalFavorited 和 点赞数量FavoriteCount
func (f *FavoriteService) GetFavoriteCntByUId(ctx context.Context, req *pb.GetFavoriteCntByUIdRequest) (*pb.GetFavoriteCntByUIdResponse, error) {

	TotalFavorited, FavoriteCount, err := f.fc.GetFavoriteCntByUId(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetFavoriteCntByUIdResponse{
		TotalFavorited: TotalFavorited,
		FavoriteCount:  FavoriteCount,
	}, nil
}
