package biz

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"video/internal/pkg/errno"

	vpb "video/api/video/v1"
	"video/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/segmentio/kafka-go"
	"golang.org/x/sync/errgroup"
)

//go:generate mockgen -destination=../mocks/mrepo/video.go -package=mrepo . VideoRepo
type VideoRepo interface {
	Save(ctx context.Context, video *model.Video) error
	PublishKafka(ctx context.Context, videoKafkaMessage *model.VideoKafkaMessage)

	GetvideoByVId(ctx context.Context, videoId int64) (*model.Video, error)
	GetVideoListByLatestTime(ctx context.Context, latestTime time.Time) ([]*model.Video, time.Time, error)

	// 获取该用户发布视频列表
	GetVideoListByAuthorId(ctx context.Context, authorId int64) ([]*model.Video, error)

	GetAuthorInfoById(ctx context.Context, authorId int64) (*vpb.User, error)
	// GetVideoListByFavoriteIdList(ctx context.Context, FavoriteList []int64) (*[]model.model.Video, error)

	//  获取视频喜欢数目
	GetFavoriteCntByVId(ctx context.Context, videoId int64) (int64, error)
	//获取视频评论数量
	GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error)
	//获取用户对视频是否喜欢
	GetIsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error)

	//redis 获取30视频id列表
	RZSetVideoIds(c context.Context, key string, latestTime int64) (int64, []int64, error)
	//30个视频id列表存入redis
	RZSetSaveVIds(ctx context.Context, value []int64, scores []float64) error

	RGetVideoInfo(ctx context.Context, videoId int64) (*model.Video, error)
	RSaveVideoInfo(ctx context.Context, video *model.Video) error

	//redis 获取该用户发布视频ID列表
	RPublishVidsByAuthorId(ctx context.Context, authorId int64) ([]int64, error)
	//redis 该用户发布视频ID列表存入redis
	RSavePublishVids(ctx context.Context, authorId int64, videoIds []int64) error
	// GetRespVideo(ctx context.Context, videoList *[]model.model.Video, userId int64) (*[]VideoList, error)
	// GetUserCountByAid

	// GetAuthorInfoByredis(c context.Context, userId int64, authorId int64) (*model.User, error)
	// GetAuthorInfoBymysql(c context.Context, userId int64, authorId int64) (*model.User, error)

	// GetIsFavoriteByRedis(c context.Context, videoId int64, userId int64) (bool, error)
	// GetIsFollowByRedis(c context.Context, userId int64, authorId int64) (bool, error)

}

type VideoUsecase struct {
	repo   VideoRepo
	log    *log.Helper
	reader *kafka.Reader
}

func NewVideoUsecase(repo VideoRepo, logger log.Logger) *VideoUsecase {
	return &VideoUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (v *VideoUsecase) Feed(ctx context.Context, latestTime time.Time, userId int64) ([]*vpb.Video, int64, error) {
	start := time.Now()
	// 获取视频列表
	// 从Redis中获取30视频id 列表
	key := "videoAll"
	earliestTime, videoIds, err := v.repo.RZSetVideoIds(ctx, key, latestTime.Unix())
	if errors.Is(err, errno.ErrRedisVIdNotFound) {
		// 从数据库中获取30视频Id列表
		video30s := make([]*model.Video, 0, 30)
		scores := make([]float64, 30)
		var nextTime time.Time

		video30s, nextTime, err = v.repo.GetVideoListByLatestTime(ctx, latestTime)
		if err != nil {
			v.log.Error("biz.Feed/GetVideoListByLatestTime", err)
			return nil, time.Now().Unix(), nil
		}
		earliestTime = nextTime.Unix()

		for k, v := range video30s {
			videoIds[k] = v.Id
			scores[k] = float64(v.CreatedAt.Unix())
		}

		//存入redis
		err = v.repo.RZSetSaveVIds(ctx, videoIds, scores)
		if err != nil {
			v.log.Error("biz.Feed/RZSetSaveVIds", err)
			return nil, time.Now().Unix(), nil
		}

	} else if err != nil {
		v.log.Error("biz.Feed/RZSetVideoIds", err)
		return nil, time.Now().Unix(), err
	}

	// 根据 video_id 从缓存中查询 video_info
	videos := make([]*model.Video, len(videoIds))
	for i, videoId := range videoIds {
		RGetVideoInfo, err := v.repo.RGetVideoInfo(ctx, videoId)
		if errors.Is(err, errno.ErrRedisVInfoNotFound) {
			// 从数据库中获取视频信息
			RGetVideoInfo, err = v.repo.GetvideoByVId(ctx, videoId)
			if err != nil {
				v.log.Error("biz.Feed/GetvideoByVId", err)
				return nil, time.Now().Unix(), err
			}
			// 存入redis
			err = v.repo.RSaveVideoInfo(ctx, RGetVideoInfo)
			if err != nil {
				v.log.Error("biz.Feed/RSaveVideoInfo", err)
				return nil, time.Now().Unix(), err
			}

		} else if err != nil {
			v.log.Error("biz.Feed/RGetVideoInfo", err)
			return nil, time.Now().Unix(), err

		}
		videos[i] = RGetVideoInfo
	}

	respVideo, err := v.GetRespVideo(ctx, videos, userId)
	if err != nil {
		v.log.Error("biz.Feed/GetRespVideo", err)
		return nil, time.Now().Unix(), err
	}

	v.log.Infof("biz.Feed success , biz.Feed耗时=%v", time.Since(start))
	return respVideo, earliestTime, nil
}

func (v *VideoUsecase) Publish(ctx context.Context, title string, videoData *[]byte, userId int64) error {
	start := time.Now()
	// 检查是否为视频文件
	if !filetype.IsVideo(*videoData) {
		return fmt.Errorf("file is not a video")
	}

	//currentDir, _ := os.Getwd()
	//v.log.Debug(" Publish工作目录:", currentDir)

	// 设置存储路径和文件名
	storagePath := "../../store/video"       // 设置存储目录
	filename := uuid.New().String() + ".mp4" // 设置文件名或使用其他方法生成唯一的文件名

	// 创建存储目录（如果不存在）
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		v.log.Error("biz.Publish/MkdirAll", err)
		return err
	}

	// 拼接完整的文件路径
	filePath := filepath.Join(storagePath, filename)

	// 写入视频数据到文件
	if err := os.WriteFile(filePath, *videoData, 0644); err != nil {
		v.log.Error("biz.Publish/WriteFile", err)
		return err
	}

	videoKafkaMessage := &model.VideoKafkaMessage{}
	videoKafkaMessage.AuthorId = userId
	videoKafkaMessage.Title = title
	videoKafkaMessage.VideoPath = filePath
	videoKafkaMessage.VideoFileName = filename

	//v.log.Debug("videoKafkaMessage:", videoKafkaMessage)
	//存储视频信息到kafka
	go v.repo.PublishKafka(ctx, videoKafkaMessage)

	v.log.Infof("biz.Publish success , biz.Publish耗时=%v", time.Since(start))
	return nil
}

func (v *VideoUsecase) PublishList(ctx context.Context, userId int64) ([]*vpb.Video, error) {
	start := time.Now()

	videoIds, err := v.repo.RPublishVidsByAuthorId(ctx, userId)
	if errors.Is(err, errno.ErrRedisPublishVidsNotFound) {
		//  获取该用户发布视频列表
		videos, err := v.repo.GetVideoListByAuthorId(ctx, userId)
		if err != nil {
			v.log.Error("biz.PublishList/GetVideoListByAuthorId", err)
			return nil, err
		}
		for _, video := range videos {
			videoIds = append(videoIds, video.Id)

		}
		//存入redis
		err = v.repo.RSavePublishVids(ctx, userId, videoIds)
		if err != nil {
			v.log.Error("biz.PublishList/RSavePublishVids", err)
			return nil, err

		}

	} else if err != nil {
		v.log.Error("biz.PublishList/RPublishVidsByAuthorId", err)
		return nil, err

	}

	// 根据 video_id 从缓存中查询 video_info
	videoList := make([]*model.Video, len(videoIds))
	for i, videoId := range videoIds {
		RGetVideoInfo, err := v.repo.RGetVideoInfo(ctx, videoId)
		if errors.Is(err, errno.ErrRedisVInfoNotFound) {
			// 从数据库中获取视频信息
			RGetVideoInfo, err = v.repo.GetvideoByVId(ctx, videoId)
			if err != nil {
				v.log.Error("biz.PublishList/GetvideoByVId", err)
				return nil, err
			}
			// 存入redis
			err = v.repo.RSaveVideoInfo(ctx, RGetVideoInfo)
			if err != nil {
				v.log.Error("biz.PublishList/RSaveVideoInfo", err)
				return nil, err
			}

		} else if err != nil {
			v.log.Error("biz.PublishList/RGetVideoInfo", err)
			return nil, err

		}
		videoList[i] = RGetVideoInfo
	}

	respVideo, err := v.GetRespVideo(ctx, videoList, userId)
	if err != nil {
		v.log.Error("biz.PublishList/GetRespVideo ", err)
		return nil, err
	}

	v.log.Infof("biz.PublishList success , biz.PublishList耗时=%v", time.Since(start))
	return respVideo, nil

}

func (v *VideoUsecase) WorkCnt(ctx context.Context, userId int64) (int64, error) {
	start := time.Now()
	videoIds, err := v.repo.RPublishVidsByAuthorId(ctx, userId)
	if errors.Is(err, errno.ErrRedisPublishVidsNotFound) {
		//  获取该用户发布视频列表
		videos, err := v.repo.GetVideoListByAuthorId(ctx, userId)
		if err != nil {
			v.log.Error("biz.WorkCnt/GetVideoListByAuthorId", err)
			return 0, err
		}
		for _, video := range videos {
			videoIds = append(videoIds, video.Id)
		}
		//存入redis
		err = v.repo.RSavePublishVids(ctx, userId, videoIds)
		if err != nil {
			v.log.Error("biz.WorkCnt/RSavePublishVids:", err)
			return 0, err

		}

	} else if err != nil {
		v.log.Error("biz.WorkCnt/RPublishVidsByAuthorId", err)
		return 0, err

	}
	v.log.Infof("biz.WorkCnt success , biz.WorkCnt耗时=%v", time.Since(start))
	return int64(len(videoIds)), nil
}

func (v *VideoUsecase) FavoriteListByVId(ctx context.Context, VideoIdList []int64) ([]*vpb.Video, error) {
	start := time.Now()
	// 根据 video_id 从缓存中查询 video_info
	videoList := make([]*model.Video, len(VideoIdList))
	for i, videoId := range VideoIdList {
		RGetVideoInfo, err := v.repo.RGetVideoInfo(ctx, videoId)
		if errors.Is(err, errno.ErrRedisVInfoNotFound) {
			// 从数据库中获取视频信息
			RGetVideoInfo, err = v.repo.GetvideoByVId(ctx, videoId)
			if err != nil {
				v.log.Debug("biz.FavoriteListByVId/GetvideoByVId", err)
				return nil, err
			}
			// 存入redis
			err = v.repo.RSaveVideoInfo(ctx, RGetVideoInfo)
			if err != nil {
				v.log.Error("biz.FavoriteListByVId/RSaveVideoInfo", err)
				return nil, err
			}

		} else if err != nil {
			v.log.Error("biz.FavoriteListByVId/RGetVideoInfo", err)
			return nil, err

		}
		videoList[i] = RGetVideoInfo
	}

	respVideo, err := v.GetRespVideo(ctx, videoList, 0)
	if err != nil {
		v.log.Error("biz.FavoriteListByVId/GetRespVideo", err)
		return nil, err
	}

	v.log.Infof("biz.FavoriteListByVId success , biz.FavoriteListByVId耗时=%v", time.Since(start))
	return respVideo, nil

}

func (v *VideoUsecase) PublishVidsByAId(ctx context.Context, authorId int64) ([]int64, error) {
	start := time.Now()
	videoIds, err := v.repo.RPublishVidsByAuthorId(ctx, authorId)
	if errors.Is(err, errno.ErrRedisPublishVidsNotFound) {
		//  获取该用户发布视频列表
		videos, err := v.repo.GetVideoListByAuthorId(ctx, authorId)
		if err != nil {
			v.log.Error("biz.PublishVidsByAId/GetVideoListByAuthorId", err)
			return nil, err
		}
		for _, video := range videos {
			videoIds = append(videoIds, video.Id)

		}
		//存入redis
		err = v.repo.RSavePublishVids(ctx, authorId, videoIds)
		if err != nil {
			v.log.Error("biz.PublishVidsByAId/RSavePublishVids", err)
			return nil, err

		}

	} else if err != nil {
		v.log.Error("biz.PublishVidsByAId/RPublishVidsByAuthorId", err)
		return nil, err

	}
	v.log.Infof("biz.PublishVidsByAId success , biz.PublishVidsByAId耗时=%v", time.Since(start))
	return videoIds, nil
}

func (v *VideoUsecase) GetAIdByVId(ctx context.Context, videoId int64) (int64, error) {
	start := time.Now()
	video, err := v.repo.GetvideoByVId(ctx, videoId)
	if err != nil {
		v.log.Error("biz.GetAIdByVId/GetvideoByVId", err)
		return 0, err
	}
	v.log.Infof("biz.GetAIdByVId success , biz.GetAIdByVId耗时=%v", time.Since(start))
	return video.AuthorId, nil
}

// 注意userId和authorId
func (v *VideoUsecase) GetRespVideo(ctx context.Context, videoList []*model.Video, userId int64) ([]*vpb.Video, error) {
	start := time.Now()
	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 对每个视频启动一个 goroutine 并在内部并发获取各项数据
	for _, video := range videoList {
		video := video // 捕获循环变量
		eg.Go(func() error {
			// 使用子 errgroup 并发获取各数据项
			var (
				author        *vpb.User
				favoriteCount int64
				commentCount  int64
				isFavorite    bool
			)
			eg2, ctx2 := errgroup.WithContext(ctx)

			// 获取作者信息
			eg2.Go(func() error {
				a, err := v.repo.GetAuthorInfoById(ctx2, video.AuthorId)
				if err != nil {
					v.log.Errorw("Failed to get author info", "videoId", video.Id, "err", err)
					return err
				}
				author = a
				return nil
			})

			// 获取视频喜欢数目
			eg2.Go(func() error {
				fc, err := v.repo.GetFavoriteCntByVId(ctx2, video.AuthorId)
				if err != nil {
					v.log.Errorw("Failed to get favorite count", "videoId", video.Id, "err", err)
					return err
				}
				favoriteCount = fc
				return nil
			})

			// 获取视频评论数量
			eg2.Go(func() error {
				cc, err := v.repo.GetCommentCntByVId(ctx2, video.AuthorId)
				if err != nil {
					v.log.Errorw("Failed to get comment count", "videoId", video.Id, "err", err)
					return err
				}
				commentCount = cc
				return nil
			})

			// 获取用户对视频是否喜欢
			eg2.Go(func() error {
				fav, err := v.repo.GetIsFavorite(ctx2, video.Id, userId)
				if err != nil {
					v.log.Errorw("Failed to get isFavorite", "videoId", video.Id, "err", err)
					return err
				}
				isFavorite = fav
				return nil
			})

			// 等待所有子任务完成
			if err := eg2.Wait(); err != nil {
				return err
			}

			// 构造最终的响应数据，一次性更新
			finalVideo := &vpb.Video{
				Id:            video.Id,
				PlayUrl:       video.PlayUrl,
				CoverUrl:      video.CoverUrl,
				Title:         video.Title,
				Author:        author,
				FavoriteCount: favoriteCount,
				CommentCount:  commentCount,
				IsFavorite:    isFavorite,
			}
			m.Store(video.Id, finalVideo)
			return nil
		})
	}

	// 等待所有视频任务完成
	if err := eg.Wait(); err != nil {
		v.log.Errorw("Failed to wait for all video tasks", "err", err)
		return nil, err
	}

	// 从 sync.Map 中一次性构造响应列表
	resp := make([]*vpb.Video, 0, len(videoList))
	for _, val := range videoList {
		videoData, ok := m.Load(val.Id)
		if !ok {
			v.log.Error("Failed to load video from sync.Map", "video_id", val.Id)
			return nil, errors.New("failed to load video from sync.Map")
		}
		resp = append(resp, videoData.(*vpb.Video))
	}

	v.log.Infof("biz.GetRespVideo success, 耗时=%v", time.Since(start))
	return resp, nil
}
