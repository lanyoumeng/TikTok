package biz

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
	"video/internal/pkg/errno"

	vpb "video/api/video/v1"
	"video/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
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

	return respVideo, earliestTime, nil
}

func (v *VideoUsecase) Publish(ctx context.Context, title string, videoData *[]byte, userId int64) error {

	currentDir, _ := os.Getwd()
	v.log.Debug(" Publish工作目录:", currentDir)

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
	v.repo.PublishKafka(ctx, videoKafkaMessage)

	return nil
}

func (v *VideoUsecase) PublishList(ctx context.Context, userId int64) ([]*vpb.Video, error) {
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
		v.log.Error("biz.PublishList/GetRespVideo", err)
		return nil, err
	}
	return respVideo, nil

}

func (v *VideoUsecase) WorkCnt(ctx context.Context, userId int64) (int64, error) {
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
	return int64(len(videoIds)), nil
}

func (v *VideoUsecase) FavoriteListByVId(ctx context.Context, VideoIdList []int64) ([]*vpb.Video, error) {

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
	return respVideo, nil

}

func (v *VideoUsecase) PublishVidsByAId(ctx context.Context, authorId int64) ([]int64, error) {
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
	return videoIds, nil
}

func (v *VideoUsecase) GetAIdByVId(ctx context.Context, videoId int64) (int64, error) {
	video, err := v.repo.GetvideoByVId(ctx, videoId)
	if err != nil {
		v.log.Error("biz.GetAIdByVId/GetvideoByVId", err)
		return 0, err
	}
	return video.AuthorId, nil
}

// 注意userId和authorId
func (v *VideoUsecase) GetRespVideo(ctx context.Context, videoList []*model.Video, userId int64) ([]*vpb.Video, error) {

	var m sync.Map
	eg, ctx := errgroup.WithContext(ctx)

	// 使用 goroutine 提高接口性能
	// 并发RPC调用获取author信息、favoriteCount、commentCount和isFavorite
	for _, video := range videoList {

		//获取作者信息  video 的redis没有缓存
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				author, err := v.repo.GetAuthorInfoById(ctx, video.AuthorId)
				if err != nil {
					v.log.Errorw("Failed to list posts", "err", err)
					return err
				}

				// 读取数据
				value, exists := m.Load(video.Id)
				if exists {
					// 数据存在，进行修改并再次写入
					if existVideo, ok := value.(*vpb.Video); ok {
						// 在现有数据的基础上修改
						existVideo.Id = video.Id
						existVideo.PlayUrl = video.PlayUrl
						existVideo.CoverUrl = video.CoverUrl
						existVideo.Title = video.Title
						existVideo.Author = author

						// 将修改后的数据重新存储到 sync.Map 中
						m.Store(video.Id, existVideo)
					} else {
						// 类型断言失败，记录错误并处理
						v.log.Errorw("Failed to assert type for video ID", "id", video.Id, "value", value)
						return errors.New("类型断言失败")
					}
				} else {
					// 数据不存在，进行初始化并存储
					m.Store(video.Id, &vpb.Video{
						Id:       video.Id,
						PlayUrl:  video.PlayUrl,
						CoverUrl: video.CoverUrl,
						Title:    video.Title,
						Author:   author,
					})
				}

				return nil
			}
		})

		//获取视频喜欢数目
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				favorCount, err := v.repo.GetFavoriteCntByVId(ctx, video.AuthorId)
				if err != nil {
					v.log.Errorw("v.repo.GetFavoriteCntByVId(ctx, video.AuthorId)", "err", err)
					return err
				}
				// 读取数据
				value, exists := m.Load(video.Id)
				if exists {
					// 数据存在，进行修改并再次写入
					if existVideo, ok := value.(*vpb.Video); ok {
						// 在现有数据的基础上修改
						existVideo.FavoriteCount = favorCount
						// 将修改后的数据重新存储到 sync.Map 中
						m.Store(video.Id, existVideo)
					} else {
						// 类型断言失败，记录错误并处理
						v.log.Errorw("Failed to assert type for video ID", "id", video.Id, "value", value)
						return errors.New("类型断言失败")
					}
				} else {
					// 数据不存在，进行初始化并存储
					m.Store(video.Id, &vpb.Video{
						FavoriteCount: favorCount,
					})
				}

				return nil
			}
		})
		//获取视频评论数量
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				commentCount, err := v.repo.GetCommentCntByVId(ctx, video.AuthorId)
				if err != nil {
					v.log.Errorw("v.repo.GetFavoriteCntByVId(ctx, video.AuthorId)", "err", err)
					return err
				}
				// 读取数据
				value, exists := m.Load(video.Id)
				if exists {
					// 数据存在，进行修改并再次写入
					if existVideo, ok := value.(*vpb.Video); ok {
						// 在现有数据的基础上修改
						existVideo.CommentCount = commentCount
						// 将修改后的数据重新存储到 sync.Map 中
						m.Store(video.Id, existVideo)
					} else {
						// 类型断言失败，记录错误并处理
						v.log.Errorw("Failed to assert type for video ID", "id", video.Id, "value", value)
						return errors.New("类型断言失败")
					}
				} else {
					// 数据不存在，进行初始化并存储
					m.Store(video.Id, &vpb.Video{
						CommentCount: commentCount,
					})
				}

				return nil
			}
		})
		//获取用户对视频是否喜欢
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				isFavorite, err := v.repo.GetIsFavorite(ctx, video.Id, userId)
				if err != nil {
					v.log.Errorw("v.repo.GetIsFavorite(ctx, video.AuthorId, userId)", "err", err)
					return err
				}
				// 读取数据
				value, exists := m.Load(video.Id)
				if exists {
					// 数据存在，进行修改并再次写入
					if existVideo, ok := value.(*vpb.Video); ok {
						// 在现有数据的基础上修改
						existVideo.IsFavorite = isFavorite
						// 将修改后的数据重新存储到 sync.Map 中
						m.Store(video.Id, existVideo)
					} else {
						// 类型断言失败，记录错误并处理
						v.log.Errorw("Failed to assert type for video ID", "id", video.Id, "value", value)
						return errors.New("类型断言失败")
					}
				} else {
					// 数据不存在，进行初始化并存储
					m.Store(video.Id, &vpb.Video{
						IsFavorite: isFavorite,
					})
				}
				return nil
			}
		})

	}

	if err := eg.Wait(); err != nil {
		v.log.Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}
	// var resp []*pb.model.Video
	resp := make([]*vpb.Video, 0, len(videoList))

	for _, v := range videoList {
		user, _ := m.Load(v.Id)
		resp = append(resp, user.(*vpb.Video))
	}

	return resp, nil
}
