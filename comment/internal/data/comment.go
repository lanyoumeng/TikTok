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
	start := time.Now()
	//评论发布日期，格式 mm-dd
	comment.CreateDate = time.Now().Format("01-02")
	err := r.data.db.Save(comment).Error
	if err != nil {
		r.log.Errorf("SaveComment-err: %v", err)
		return nil, err
	}

	r.log.Infof("SaveComment success , commentId=%v , SaveComment耗时=%v", comment.Id, time.Since(start))
	return comment, nil
}

func (r *commentRepo) DelComment(ctx context.Context, commentId int64) error {
	start := time.Now()
	err := r.data.db.Where("id = ?", commentId).Delete(&model.Comment{}).Error
	if err != nil {
		r.log.Errorf("DelComment-err: %v", err)
		return err
	}

	r.log.Infof("DelComment success , commentId=%v , DelComment耗时=%v", commentId, time.Since(start))
	return nil
}

// 根据userId查询用户信息
func (r *commentRepo) GetUserinfoByUId(ctx context.Context, userId int64) (*pb.User, error) {
	start := time.Now()
	resp, err := r.data.userc.UserInfo(context.Background(), &userV1.UserRequest{UserId: strconv.FormatInt(userId, 10)})
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

	r.log.Infof("GetUserinfoByUId success , 耗时=%v", time.Since(start))
	return userInfo, nil
}

// 根据videoId查询作者Id
func (r *commentRepo) GetAuthorIdByVId(ctx context.Context, videoId int64) (int64, error) {
	start := time.Now()
	resp, err := r.data.videoc.GetAIdByVId(context.Background(), &videoV1.GetAIdByVIdReq{VideoId: videoId})
	if err != nil {
		r.log.Errorf("Error getting author info: %v", err)
		return 0, err

	}

	r.log.Infof("GetAuthorIdByVId success , 耗时=%v", time.Since(start))
	return resp.AuthorId, nil

}

// 根据userId,authorId查询用户是否关注作者
func (r *commentRepo) GetFollowByUIdAId(ctx context.Context, userId, authorId int64) (bool, error) {
	start := time.Now()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*3000)
	//defer cancel()

	resp, err := r.data.relationc.IsFollow(context.Background(), &relationV1.IsFollowRequest{UserId: userId, AuthorId: authorId})
	if err != nil {
		r.log.Errorf("Error getting follow info: %v", err)
		return false, err
	}

	r.log.Infof("GetFollowByUIdAId success , 耗时=%v", time.Since(start))
	return resp.IsFollow, nil
}

// 评论列表
func (r *commentRepo) CommentList(ctx context.Context, videoId int64) ([]*model.Comment, error) {
	start := time.Now()
	//zset
	//comment::video_id
	//score: 评论发布时间戳 member: id+user_id+content
	//从缓存获取
	commentList := make([]*model.Comment, 0)

	//key
	key := "comment::" + strconv.FormatInt(videoId, 10)

	//r.log.Infof("key: %v", key)
	//从缓存获取
	commentSet, err := r.data.rdb.ZRangeWithScores(context.Background(), key, 0, -1).Result()
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

			//CreateDate是评论发布日期，格式 mm-dd 不能直接转换为时间戳，用UpdatedAt代替
			//score, err := strconv.ParseFloat(comment.CreateDate, 64)
			//if err != nil {
			//	r.log.Fatalf("strconv.ParseFloat failed: %v", err)
			//}
			score := float64(comment.UpdatedAt.Unix())

			//存入缓存
			err = r.data.rdb.ZAdd(context.Background(), key, &redis.Z{Score: score, Member: member}).Err()
			if err != nil {
				r.log.Errorf("add cache failed: %v", err)
			}

		}

		//过期时间
		err = r.data.rdb.Expire(context.Background(), "comment::"+strconv.FormatInt(videoId, 10), tool.GetRandomExpireTime()).Err()
		if err != nil {
			r.log.Errorf("set redis failed: %v", err)
			return nil, err
		}

		r.log.Infof("mysql----len(commentList):%v", len(commentList))
		//格式对 直接返回
		return commentList, nil

	} else if err != nil {
		r.log.Errorf("get cache failed: %v", err)
		return nil, err
	}

	//可以从缓存获取    数据转换为model.Comment
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
	r.log.Infof("redis----len(commentList):%v", len(commentList))

	r.log.Infof("CommentList success , 耗时=%v", time.Since(start))
	return commentList, nil
}

// 获取评论数
func (r *commentRepo) GetCommentCntByVId(ctx context.Context, videoId int64) (int64, error) {
	start := time.Now()
	//zset
	//comment::video_id
	//score: 评论发布时间戳 member: id+user_id+content
	//从缓存获取
	count, err := r.data.rdb.ZCard(context.Background(), "comment::"+strconv.FormatInt(videoId, 10)).Result()
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

	r.log.Infof("GetCommentCntByVId success , 耗时=%v", time.Since(start))
	return count, nil
}
