package data

import (
	"context"
	"sync"
	"time"
	favoriteV1 "user/api/favorite/v1"
	relationV1 "user/api/relation/v1"
	"user/internal/biz"
	"user/internal/pkg/errno"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"

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
	start := time.Now()
	res := r.data.db.Model(&model.User{}).Create(&u)
	if res.Error != nil {
		r.log.Error("CreateUser-err:", res.Error)
		return 0, res.Error
	}

	r.log.Infof("CreateUser success , userId=%v , 耗时=%v", u.Id, time.Since(start))
	return u.Id, nil
}

func (r *userRepo) UserByName(ctx context.Context, name string) (*model.User, error) {
	start := time.Now()
	user := &model.User{}
	err := r.data.db.Model(&model.User{}).Where("name = ?", name).First(&user).Error
	if err != nil {
		r.log.Error("UserByName-err:", err)
		return nil, err
	}
	r.log.Infof("UserByName success , 耗时=%v", time.Since(start))
	return user, nil
}

// GetUserById .
func (r *userRepo) GetUserById(ctx context.Context, Id int64) (*model.User, error) {
	start := time.Now()
	user := &model.User{}
	if err := r.data.db.Model(&model.User{}).Where("id = ?", Id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Error("GetUserById-err:", errno.ErrUserNotFound)
			return nil, errno.ErrUserNotFound
		}
		r.log.Error("GetUserById-err:", err)
		return nil, err
	}
	//如果用户不存在，报错
	if user.Id == 0 {
		r.log.Error("user.Id == 0, GetUserById-err:", errno.ErrUserNotFound)
		return nil, errno.ErrUserNotFound
	}

	r.log.Infof("GetUserById success , 耗时=%v", time.Since(start))
	return user, nil
}

// UpdateUser .
func (r *userRepo) UpdateUser(ctx context.Context, user *model.User) error {
	start := time.Now()
	if err := r.data.db.Model(&model.User{}).Save(&user).Error; err != nil {
		r.log.Error("UpdateUser-err:", err)
		return err
	}

	r.log.Infof("UpdateUser success , 耗时=%v", time.Since(start))
	return nil
}

// 调用rpc服务 获取用户计数信息
func (r *userRepo) GetCountById(ctx context.Context, id int64) (*model.UserCount, error) {
	start := time.Now()
	count := &model.UserCount{}
	count.Id = id

	var wg sync.WaitGroup

	wg.Add(3)

	//FollowCount  FollowerCnt
	go func() {
		defer wg.Done()
		start := time.Now()
		followCntRes, err := r.data.relationc.FollowCnt(context.Background(), &relationV1.FollowCntRequest{UserId: id})
		if err != nil {
			r.log.Errorf("data.GetCountById/FollowCnt-err:%v\n", err)
			return
		}
		count.FollowCount = followCntRes.FollowCnt
		count.FollowerCount = followCntRes.FollowerCnt
		r.log.Infof("data.GetCountById/FollowCnt success , userId=%v , 耗时=%v", id, time.Since(start))
	}()

	//WorkCount
	go func() {
		defer wg.Done()
		start := time.Now()
		workCntRes, err := r.data.videoc.WorkCnt(context.Background(), &videov1.WorkCntRequest{UserId: id})
		if err != nil {
			r.log.Errorf("data.GetCountById/WorkCnt-err:%v\n", err)
			return
		}
		count.WorkCount = workCntRes.WorkCount
		r.log.Infof("data.GetCountById/WorkCnt success , userId=%v , 耗时=%v", id, time.Since(start))
	}()

	//FavoriteCount  TotalFavorited
	go func() {
		defer wg.Done()
		start := time.Now()
		favoriteCntRes, err := r.data.favc.GetFavoriteCntByUId(context.Background(), &favoriteV1.GetFavoriteCntByUIdRequest{UserId: id})
		if err != nil {
			r.log.Errorf("data.GetCountById/GetFavoriteCntByUId-err:%v\n", err)
			return
		}
		count.FavoriteCount = favoriteCntRes.FavoriteCount
		count.TotalFavorited = favoriteCntRes.TotalFavorited
		r.log.Infof("data.GetCountById/GetFavoriteCntByUId success , userId=%v , 耗时=%v", id, time.Since(start))
	}()

	wg.Wait()

	r.log.Infof("GetCountById success , userId=%v , 耗时=%v", id, time.Since(start))
	return count, nil

}
