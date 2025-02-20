package biz

import (
	"context"
	"time"

	vpb "favorite/api/video/v1"

	"github.com/go-kratos/kratos/v2/log"
)

//go:generate mockgen -destination=../mocks/mrepo/favorite.go -package=mrepo . FavoriteRepo
type FavoriteRepo interface {
	GetAuthorIdByVideoId(ctx context.Context, videoId int64) (int64, error)
	//点赞还是取消
	FavoriteAction(ctx context.Context, authorId int64, videoId int64, userId int64, action_type int64) error
	FavoriteList(ctx context.Context, userId int64) ([]*vpb.Video, error)
	IsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error)
	// PublishKafka(ctx context.Context, favoriteKafkaMessage *model.FavoriteKafkaMessage)

	// GetlikeIdListByUserId(ctx context.Context, userId int64) ([]int64, error)
	GetFavoriteCntByVId(ctx context.Context, videoId int64) (int64, error)
	// GetIslike(ctx context.Context, videoId int64, userId int64) (bool, error)

	//获取 用户的 获赞数TotalFavorited 和 点赞数量FavoriteCount
	GetFavoriteCntByUId(ctx context.Context, userId int64) (int64, int64, error)
}

type FavoriteUsecase struct {
	repo FavoriteRepo
	log  *log.Helper
}

func NewFavoriteUsecase(repo FavoriteRepo, logger log.Logger) *FavoriteUsecase {
	return &FavoriteUsecase{repo: repo, log: log.NewHelper(logger)}
}
func (f *FavoriteUsecase) GetAuthorIdByVideoId(ctx context.Context, videoId int64) (int64, error) {
	start := time.Now()
	authorId, err := f.repo.GetAuthorIdByVideoId(ctx, videoId)
	if err != nil {
		f.log.Error("GetAuthorIdByVideoId err:", err)
		return 0, err
	}
	f.log.Infof("GetAuthorIdByVideoId success , GetAuthorIdByVideoId耗时=%v", time.Since(start))
	return authorId, nil
}

func (f *FavoriteUsecase) Favorite(ctx context.Context, authorId int64, videoId int64, userId int64, actionType int64) error {
	start := time.Now()
	//1点赞 2取消点赞
	if err := f.repo.FavoriteAction(ctx, authorId, videoId, userId, actionType); err != nil {
		f.log.Error("FavoriteAction err:", err)
		return err
	}

	f.log.Infof("Favorite success , Favorite耗时=%v", time.Since(start))
	return nil
}

func (f *FavoriteUsecase) FavoriteList(ctx context.Context, userId int64) ([]*vpb.Video, error) {
	start := time.Now()
	// var videoList []*vpb.Video
	videoList, err := f.repo.FavoriteList(ctx, userId)
	if err != nil {
		f.log.Error("FavoriteList err:", err)
		return nil, err
	}

	f.log.Infof("FavoriteList success , FavoriteList耗时=%v", time.Since(start))
	return videoList, nil
}

func (f *FavoriteUsecase) GetFavoriteCntByVId(ctx context.Context, videoId int64) (int64, error) {
	start := time.Now()
	favoriteCnt, err := f.repo.GetFavoriteCntByVId(ctx, videoId)
	if err != nil {
		f.log.Error("GetFavoriteCntByVId err:", err)
		return 0, err
	}
	f.log.Infof("GetFavoriteCntByVId success , GetFavoriteCntByVId耗时=%v", time.Since(start))
	return favoriteCnt, nil
}

func (f *FavoriteUsecase) GetIsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	start := time.Now()
	flag, err := f.repo.IsFavorite(ctx, videoId, userId)
	if err != nil {
		f.log.Error("GetIsFavorite err:", err)
		return false, err
	}
	f.log.Infof("GetIsFavorite success , GetIsFavorite耗时=%v", time.Since(start))
	return flag, nil
}

// 获取 用户的 获赞数favedCnt 和 点赞数量favCnt
func (f *FavoriteUsecase) GetFavoriteCntByUId(ctx context.Context, userId int64) (int64, int64, error) {
	start := time.Now()
	favedCnt, favCnt, err := f.repo.GetFavoriteCntByUId(ctx, userId)
	if err != nil {
		f.log.Error("GetFavoriteCntByUId err:", err)
		return 0, 0, err
	}

	f.log.Infof("GetFavoriteCntByUId success , GetFavoriteCntByUId耗时=%v", time.Since(start))
	return favedCnt, favCnt, nil
}
