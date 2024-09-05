package data

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	cpb "video/api/comment/v1"
	fpb "video/api/favorite/v1"
	upb "video/api/user/v1"
	vpb "video/api/video/v1"
	"video/internal/biz"
	"video/internal/pkg/model"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"github.com/segmentio/kafka-go"
)

type videoRepo struct {
	data *Data
	log  *log.Helper
}

var _ biz.VideoRepo = (*videoRepo)(nil)

// New videoRepo.
func NewBizVideoRepo(data *Data, logger log.Logger) biz.VideoRepo {
	return &videoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "biz.repo/video")),
	}
}

func (v *videoRepo) Save(ctx context.Context, video *model.Video) error {

	if err := v.data.db.Model(&model.Video{}).Save(&video).Error; err != nil {
		v.log.Error("Error creating user1:", err)
	}
	return nil
}

// 发布视频上传信息·到消息队列
func (v *videoRepo) PublishKafka(ctx context.Context, videoKafkaMessage *model.VideoKafkaMessage) {
	//videoByte, _ := sonic.Marshal(videoKafkaMessage)

	//v.log.Debug("data/PublishKafka:", videoKafkaMessage)
	videoByte, err := json.Marshal(videoKafkaMessage)
	if err != nil {
		v.log.Error("PublishKafka-err:", err)
		return

	}

	for i := 0; i < 3; i++ { //允许重试3次
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := v.data.kafakProducer.WriteMessages(ctx, //批量写入消息，原子操作，要么全写成`功，要么全写失败

			kafka.Message{Value: videoByte},
		); err != nil {
			// if err == kafka.LeaderNotAvailable || errors.Is(err, context.DeadlineExceeded) {
			if errors.Is(err, kafka.LeaderNotAvailable) { //首次写一个新的Topic时，会发生LeaderNotAvailable错误，重试一次就好了
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				v.log.Errorf("batch write message failed: %v", err)
			}
		} else {
			v.log.Debug("write message to kafka success")
			break //只要成功一次就不再尝试下一次了
		}

	}

}
func (v *videoRepo) GetvideoByVId(ctx context.Context, videoId int64) (*model.Video, error) {
	video := &model.Video{}
	err := v.data.db.Model(&model.Video{}).Where("id = ?", videoId).First(&video).Error
	if err != nil {
		v.log.Error("GetvideoByVId-err:", err)
		return nil, err
	}
	return video, nil
}
func (v *videoRepo) GetVideoListByLatestTime(ctx context.Context, latestTime time.Time) ([]*model.Video, time.Time, error) {
	videos := make([]*model.Video, 10)
	err := v.data.db.Model(&model.Video{}).Where("created_at < ?", latestTime).Order("created_at desc").Limit(30).Find(&videos).Error
	if err != nil {
		v.log.Error("GetVideoListByLatestTime-err:", err)
		return nil, time.Now(), err

	}

	if len(videos) == 0 {
		return videos, time.Now(), nil
	} else {
		return videos, videos[len(videos)-1].CreatedAt, nil
	}

}

// 获取该用户发布视频列表
func (v *videoRepo) GetVideoListByAuthorId(ctx context.Context, userId int64) ([]*model.Video, error) {
	videos := make([]*model.Video, 10)
	err := v.data.db.Model(&model.Video{}).Where("author_id = ?", userId).Find(&videos).Error

	if err != nil {
		v.log.Error("GetVideoListByAuthorId-err:", err)
		return nil, err
	}

	return videos, nil

}

func (v *videoRepo) GetAuthorInfoById(ctx context.Context, authorId int64) (*vpb.User, error) {
	userrepo, err := v.data.userc.UserInfo(ctx, &upb.UserRequest{UserId: strconv.FormatInt(authorId, 10)})
	if err != nil {
		v.log.Error("GetAuthorInfoById-err:", err)
		return nil, err
	}
	var user *vpb.User
	_ = copier.Copy(user, userrepo.User)

	return user, nil

}

func (v *videoRepo) GetFavoriteCntByVId(ctx context.Context, videoId int64) (int64, error) {
	favorepo, err := v.data.favc.GetFavoriteCntByVId(ctx, &fpb.GetFavoriteCntByVIdRequest{Id: videoId})
	if err != nil {
		v.log.Error("GetFavoriteCntByVId-err:", err)
		return 0, err
	}
	return favorepo.FavoriteCount, nil
}

func (v *videoRepo) GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error) {
	commentrepo, err := v.data.commentc.GetCommentCntByVId(ctx, &cpb.GetCommentCntByVIdReq{VideoId: videoId})
	if err != nil {
		v.log.Error("GetCommentCntByVId-err:", err)
		return 0, err
	}
	return commentrepo.CommentCount, nil
}

func (v *videoRepo) GetIsFavorite(ctx context.Context, videoId int64, userId int64) (bool, error) {
	favorepo, err := v.data.favc.GetIsFavorite(ctx, &fpb.GetIsFavoriteRequest{VideoId: videoId, UserId: userId})
	if err != nil {
		v.log.Error("GetIsFavorite-err:", err)
		return false, err
	}
	return favorepo.Favorite, nil
}
