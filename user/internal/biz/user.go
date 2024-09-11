package biz

import (
	"context"

	v1 "user/api/user/v1"
	"user/internal/conf"
	"user/internal/pkg/errno"
	"user/internal/pkg/model"
	"user/pkg/auth"
	"user/pkg/token"

	"github.com/asaskevich/govalidator"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../mocks/mrepo/user.go -package=mrepo . UserRepo
type UserRepo interface {
	CreateUser(context.Context, *model.User) (int64, error)
	UserByName(ctx context.Context, name string) (*model.User, error)
	GetUserById(ctx context.Context, id int64) (*model.User, error)

	UpdateUser(context.Context, *model.User) error

	//从redis获取用户
	RGetUserById(ctx context.Context, id int64) (*model.User, error)
	//保存用户到redis
	RSaveUser(ctx context.Context, user *model.User) error

	//获取计数信息
	GetCountById(context.Context, int64) (*model.UserCount, error)
	//redis获取计数信息
	RGetCountById(context.Context, int64) (*model.UserCount, error)
	//保存/更新 计数信息到redis
	RSaveCount(context.Context, *model.UserCount) error
}

type UserUsecase struct {
	repo   UserRepo
	log    *log.Helper
	JwtKey string
}

func NewUserUsecase(repo UserRepo, logger log.Logger, conf *conf.Auth) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger), JwtKey: conf.JwtKey}
}
func (uc *UserUsecase) Create(ctx context.Context, u *model.User) (int64, string, error) {

	//uc.log.Debug("Create user")

	// 查询用户是否已存在
	user, err := uc.repo.UserByName(ctx, u.Name)
	// 如果查询出错，且不是记录不存在的错误，则返回错误
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		uc.log.Error("biz.Create/!errors.Is(err, gorm.ErrRecordNotFound)-err:", err)
		return 0, "", err
	}

	if user != nil {
		uc.log.Error("biz.Create/user!=nil-err:", err)
		return 0, "", errno.ErrUserAlreadyExist
	}

	// 验证输入数据
	if _, err := govalidator.ValidateStruct(u); err != nil {
		uc.log.Error("biz.Create/govalidator.ValidateStruct(u)-err:", err)
		return 0, "", err
	}
	// 创建一个新的用户对象，并从输入的 u 复制数据
	newUser := &model.User{}
	_ = copier.Copy(newUser, u)

	// 加密密码
	newUser.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		uc.log.Error("biz.Create/auth.Encrypt(u.Password)-err:", err)
		return 0, "", err
	}

	//因为前端原因，这里Avatar用户头像和BackgroundImage 用户个人页顶部大图默认
	newUser.Avatar = "https://tiktok-example.oss-cn-beijing.aliyuncs.com/squirrel.jpg"
	newUser.BackgroundImage = "https://tiktok-example.oss-cn-beijing.aliyuncs.com/squirrel.jpg"
	newUser.Signature = "这个人很懒，什么都没留下"

	// 创建用户
	userId, err := uc.repo.CreateUser(ctx, newUser)
	if err != nil {
		uc.log.Error("biz.Create/CreateUser-err:", err)
		return 0, "", err
	}
	newUser.Id = userId

	// 创建 JWT Token
	userinfo := &token.UserClaims{}
	userinfo.UserId = userId
	userinfo.Name = newUser.Name

	token, err := token.Sign(userinfo, uc.JwtKey)
	if err != nil {
		uc.log.Error("biz.Create/Sign-err:", err)
		return 0, "", errno.ErrSignToken
	}

	return userId, token, nil
}

func (uc *UserUsecase) Login(ctx context.Context, u *model.User) (int64, string, error) {
	// 获取登录用户的所有信息
	user, err := uc.repo.UserByName(ctx, u.Name)
	if err != nil {
		uc.log.Error("biz.Login/UserByName-err:", err)
		return 0, "", err
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(user.Password, u.Password); err != nil {
		uc.log.Error("biz.Login/auth.Compare-err:", err)
		return 0, "", errno.ErrPasswordIncorrect
	}

	// 如果匹配成功，说明登录成功，签发 token 并返回
	userinfo := &token.UserClaims{}
	userinfo.UserId = user.Id
	userinfo.Name = user.Name

	token, err := token.Sign(userinfo, uc.JwtKey)
	if err != nil {
		uc.log.Error("biz.Login/Sign-err:", err)
		return 0, "", errno.ErrSignToken
	}

	return user.Id, token, nil
}

func (uc *UserUsecase) UserByName(ctx context.Context, mobile string) (*model.User, error) {
	return uc.repo.UserByName(ctx, mobile)
}

func (uc *UserUsecase) UpdateUser(ctx context.Context, user *model.User) error {
	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserUsecase) UserById(ctx context.Context, id int64) (*model.User, error) {
	return uc.repo.GetUserById(ctx, id)
}

// 获取用户所有信息
func (uc *UserUsecase) UserInfo(ctx context.Context, id int64) (*v1.User, error) {
	userinfo := &v1.User{}
	//从数据库获取用户注册信息\作品数量\点赞数量\粉丝数量\关注数量\获赞数量\ 进行整合
	user, err := uc.repo.RGetUserById(ctx, id)
	if errors.Is(err, errno.ErrUserNotFound) {
		user, err = uc.repo.GetUserById(ctx, id)
		if err != nil {
			uc.log.Error("biz.UserInfo/GetUserById-err:", err)
			return nil, err
		}
		//将用户注册信息存入redis
		err = uc.repo.RSaveUser(ctx, user)
		if err != nil {
			uc.log.Error("biz.UserInfo/RSaveUser-err:", err)
			return nil, err
		}
	} else if err != nil {
		uc.log.Error("biz.UserInfo/RGetUserById-err:", err)
		return nil, err
	}
	//uc.log.Debug("biz.UserInfo/user:", user)
	err = copier.Copy(&userinfo, user)
	if err != nil {
		uc.log.Error("biz.UserInfo/copier.Copy-err:", err)
		return nil, err
	}

	count, err := uc.repo.RGetCountById(ctx, id)
	if errors.Is(err, errno.ErrUserNotFound) {
		count, err = uc.repo.GetCountById(ctx, id)
		if err != nil {
			uc.log.Error("biz.UserInfo/GetCountById-err:", err)
			return nil, err
		}

		//将用户计数信息存入redis
		err = uc.repo.RSaveCount(ctx, count)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		uc.log.Error("biz.UserInfo/RGetCountById-err:", err)
		return nil, err
	}
	err = copier.Copy(&userinfo, count)
	if err != nil {
		uc.log.Error("biz.UserInfo/copier.Copy-err:", err)
		return nil, err
	}

	// bool is_follow 默认就行,让video服务调用 // true-已关注，false-未关注

	return userinfo, nil
}

// 无mysql表
func (uc *UserUsecase) RGetCountById(ctx context.Context, userId int64) (*model.UserCount, error) {
	count, err := uc.repo.RGetCountById(ctx, userId)
	if err == errno.ErrUserNotFound {
		count, err = uc.repo.GetCountById(ctx, userId)
		if err != nil {
			uc.log.Error("biz.RGetCountById/GetCountById-err:", err)
			return nil, err
		}

		err = uc.repo.RSaveCount(ctx, count)
		if err != nil {
			uc.log.Error("biz.RGetCountById/RSaveCount-err:", err)
			return nil, err
		}
	} else if err != nil {
		uc.log.Error("biz.RGetCountById/RGetCountById-err:", err)
		return nil, err
	}

	return count, nil
}

func (uc *UserUsecase) RSaveCount(ctx context.Context, userCount *model.UserCount) error {
	return uc.repo.RSaveCount(ctx, userCount)
}

func (uc *UserUsecase) UserInfoList(ctx context.Context, ids []int64) ([]*v1.User, error) {
	var userInfos []*v1.User
	for _, id := range ids {
		userInfo, err := uc.UserInfo(ctx, id)
		if err != nil {
			uc.log.Error("biz.UserInfoList/UserClaims-err:", err)
			return nil, err
		}
		userInfos = append(userInfos, userInfo)
	}
	return userInfos, nil
}
