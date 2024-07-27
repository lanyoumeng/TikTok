package data

import (
	pb "comment/api/comment/v1"
	relationV1 "comment/api/relation/v1"
	userV1 "comment/api/user/v1"
	videoV1 "comment/api/video/v1"
	"comment/internal/biz"
	"comment/internal/pkg/model"
	"comment/pkg/tool"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type commentRepo struct {
	data *Data
	log  *log.Helper
}

// NewCommentRepo .
func NewCommentRepo(data *Data, logger log.Logger) biz.CommentRepo {
	return &commentRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *commentRepo) SaveComment(ctx context.Context, comment *model.Comment) (*model.Comment, error) {

	//评论发布日期，格式 mm-dd
	comment.CreateDate = time.Now().Format("01-02")
	err := r.data.db.Save(comment).Error
	if err != nil {
		r.log.Errorf("SaveComment-err: %v", err)
		return nil, err
	}

	return comment, nil
}

func (r *commentRepo) DelComment(ctx context.Context, commentId int64) error {
	err := r.data.db.Delete(&model.Comment{}, commentId).Error
	if err != nil {
		r.log.Errorf("DelComment-err: %v", err)
		return err
	}
	return nil
}

// 根据userId查询用户信息
func (r *commentRepo) GetUserinfoByUId(ctx context.Context, userId int64) (*pb.User, error) {

	resp, err := r.data.userc.UserInfo(ctx, &userV1.UserRequest{UserId: strconv.FormatInt(userId, 10)})
	if err != nil {
		r.log.Errorf("Error getting user info: %v", err)
		return nil, err

	}
	userInfo := &pb.User{}
	userInfo.Id = resp.User.Id
	userInfo.Name = resp.User.Name
	userInfo.FollowCount = resp.User.FollowCount
	userInfo.FollowerCount = resp.User.FollowerCount
	userInfo.IsFollow = resp.User.IsFollow
	userInfo.Avatar = resp.User.Avatar
	userInfo.BackgroundImage = resp.User.BackgroundImage
	userInfo.Signature = resp.User.Signature
	userInfo.TotalFavorited = resp.User.TotalFavorited
	userInfo.WorkCount = resp.User.WorkCount
	userInfo.FavoriteCount = resp.User.FavoriteCount

	return userInfo, nil
}

// 根据videoId查询作者Id
func (r *commentRepo) GetAuthorIdByVId(ctx context.Context, videoId int64) (int64, error) {

	resp, err := r.data.videoc.GetAIdByVId(ctx, &videoV1.GetAIdByVIdReq{VideoId: videoId})
	if err != nil {
		r.log.Errorf("Error getting author info: %v", err)
		return 0, err

	}
	return resp.AuthorId, nil

}

// 根据userId,authorId查询用户是否关注作者
func (r *commentRepo) GetFollowByUIdAId(ctx context.Context, userId, authorId int64) (bool, error) {
	resp, err := r.data.relationc.IsFollow(ctx, &relationV1.IsFollowRequest{UserId: userId, AuthorId: authorId})
	if err != nil {
		r.log.Errorf("Error getting follow info: %v", err)
		return false, err
	}
	return resp.IsFollow, nil
}

// 评论列表
func (r *commentRepo) CommentList(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	//zset
	//comment::video_id
	//score: 评论发布时间戳 member: id+user_id+content
	//从缓存获取
	commentList := make([]*model.Comment, 10)

	//从缓存获取
	commentSet, err := r.data.rdb.ZRangeWithScores(ctx, "comment::"+strconv.FormatInt(videoId, 10), 0, -1).Result()
	if errors.Is(err, redis.Nil) || len(commentSet) == 0 {
		//从数据库获取
		err := r.data.db.Where("video_id = ?", videoId).Find(&commentList).Error
		if err != nil {
			r.log.Errorf("CommentList-err: %v", err)
			return nil, err
		}

		//存入缓存
		for _, comment := range commentList {
			//序列化
			commentJson, err := json.Marshal(comment)
			if err != nil {
				r.log.Errorf("json.Marshal failed: %v", err)
				return nil, err
			}
			//member
			member := string(commentJson)
			//score
			score, err := strconv.ParseFloat(comment.CreateDate, 64)
			if err != nil {

				r.log.Fatalf("strconv.ParseFloat failed: %v", err)
			}

			//key
			key := "comment::" + strconv.FormatInt(videoId, 10)
			//存入缓存
			err = r.data.rdb.ZAdd(ctx, key, &redis.Z{Score: score, Member: member}).Err()
			if err != nil {
				r.log.Errorf("add cache failed: %v", err)
			}

		}

		//过期时间
		err = r.data.rdb.Expire(ctx, "comment::"+strconv.FormatInt(videoId, 10), tool.GetRandomExpireTime()).Err()
		if err != nil {
			r.log.Errorf("set redis failed: %v", err)
			return nil, err
		}

		//格式对 直接返回
		return commentList, nil
	} else if err != nil {
		r.log.Errorf("get cache failed: %v", err)
		return nil, err
	}

	//从缓存获取的数据转换为model.Comment
	for _, z := range commentSet {
		comment := &model.Comment{}
		err := json.Unmarshal([]byte(z.Member.(string)), &comment)
		if err != nil {
			r.log.Errorf("json.Unmarshal failed: %v", err)
			return nil, err
		}
		// 使用 score 作为 CreateDate
		comment.CreateDate = time.Unix(int64(z.Score), 0).Format("01-02")

		commentList = append(commentList, comment)
	}
	return commentList, nil

}

// 获取评论数
func (r *commentRepo) GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error) {
	//zset
	//comment::video_id
	//score: 评论发布时间戳 member: id+user_id+content
	//从缓存获取
	count, err := r.data.rdb.ZCard(ctx, "comment::"+strconv.FormatInt(videoId, 10)).Result()
	if errors.Is(err, redis.Nil) || count == 0 {

		//从数据库获取
		err := r.data.db.Model(&model.Comment{}).Where("video_id = ?", videoId).Count(&count).Error
		if err != nil {
			r.log.Errorf("GetCommentCntByVId-err: %v", err)
			return 0, err
		}

	} else if err != nil {
		r.log.Errorf("get cache failed: %v", err)
		return 0, err
	}

	return count, nil
}
