package data

import (
	"context"
	favoriteV1 "user/api/favorite/v1"
	"user/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

	relationV1 "user/api/relation/v1"
	videov1 "user/api/video/v1"
	"user/internal/pkg/model"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo .
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "repo/user")),
	}
}

// CreateUser .
func (r *userRepo) CreateUser(ctx context.Context, u *model.User) (int64, error) {
	res := r.data.db.Model(&model.User{}).Create(&u)
	if res.Error != nil {
		return 0, res.Error
	}

	return u.Id, nil
}

// UserByName
func (r *userRepo) UserByName(ctx context.Context, name string) (*model.User, error) {

	user := &model.User{}
	err := r.data.db.Model(&model.User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return user, nil
}

// GetUserById .
func (r *userRepo) GetUserById(ctx context.Context, Id int64) (*model.User, error) {
	user := &model.User{}
	if err := r.data.db.Model(&model.User{}).Where("id = ?", Id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		return nil, err
	}

	return user, nil
}

// UpdateUser .
func (r *userRepo) UpdateUser(ctx context.Context, user *model.User) error {

	if err := r.data.db.Model(&model.User{}).Save(&user).Error; err != nil {
		return err
	}

	return nil
}

//type UserCount struct {
//	Id             int64 `json:"id" `
//	FollowCount    int64 `json:"follow_count" `
//	FollowerCount  int64 `json:"follower_count"`

//	WorkCount      int64 `json:"work_count"`  video

//	FavoriteCount  int64 `json:"favorite_count"`
//	TotalFavorited int64 `json:"total_favorited"`
//}

// 调用rpc服务 获取用户计数信息
func (r *userRepo) GetCountById(ctx context.Context, id int64) (*model.UserCount, error) {

	count := &model.UserCount{}
	count.Id = id

	//FollowCount  FollowerCnt
	followCntRes := &relationV1.FollowCntResponse{}
	followCntRes, err := r.data.relationc.FollowCnt(ctx, &relationV1.FollowCntRequest{UserId: id})
	if err != nil {
		r.log.Errorf("data.GetCountById/FollowCnt-err:%v\n", err)
		return nil, err
	}
	count.FollowCount = followCntRes.FollowCnt
	count.FollowerCount = followCntRes.FollowerCnt

	//WorkCount
	workCntRes := &videov1.WorkCntResponse{}
	workCntRes, err = r.data.videoc.WorkCnt(ctx, &videov1.WorkCntRequest{UserId: id})
	if err != nil {
		r.log.Errorf("data.GetCountById/WorkCnt-err:%v\n", err)
		return nil, err
	}
	count.WorkCount = workCntRes.WorkCount

	//FavoriteCount  TotalFavorited
	favoriteCntRes := &favoriteV1.GetFavoriteCntByUIdResponse{}
	favoriteCntRes, err = r.data.favc.GetFavoriteCntByUId(ctx, &favoriteV1.GetFavoriteCntByUIdRequest{UserId: id})
	if err != nil {
		r.log.Errorf("data.GetCountById/GetFavoriteCntByUId-err:%v\n", err)
		return nil, err
	}
	count.FavoriteCount = favoriteCntRes.FavoriteCount
	count.TotalFavorited = favoriteCntRes.TotalFavorited

	return count, nil

}
