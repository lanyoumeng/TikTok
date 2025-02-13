// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v3.12.4
// source: user/v1/user.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationUserServiceLogin = "/user.v1.UserService/Login"
const OperationUserServiceRegister = "/user.v1.UserService/Register"
const OperationUserServiceUpdateFavoriteCnt = "/user.v1.UserService/UpdateFavoriteCnt"
const OperationUserServiceUpdateFollowCnt = "/user.v1.UserService/UpdateFollowCnt"
const OperationUserServiceUpdateWorkCnt = "/user.v1.UserService/UpdateWorkCnt"
const OperationUserServiceUserInfo = "/user.v1.UserService/UserInfo"
const OperationUserServiceUserInfoList = "/user.v1.UserService/UserInfoList"

type UserServiceHTTPServer interface {
	// Login 用户登录，有 Redis
	Login(context.Context, *UserLoginRequest) (*UserLoginResponse, error)
	// Register 用户注册
	Register(context.Context, *UserRegisterRequest) (*UserRegisterResponse, error)
	// UpdateFavoriteCnt 更新点赞数量
	// FavoriteCount 和 TotalFavorited
	UpdateFavoriteCnt(context.Context, *UpdateFavoriteCntRequest) (*UpdateFavoriteCntResponse, error)
	// UpdateFollowCnt 更新关注数量
	// FollowCount 和 FollowerCount
	UpdateFollowCnt(context.Context, *UpdateFollowCntRequest) (*UpdateFollowCntResponse, error)
	// UpdateWorkCnt 更新计数
	// 更新作品数量
	UpdateWorkCnt(context.Context, *UpdateWorkCntRequest) (*UpdateWorkCntResponse, error)
	// UserInfo 获取用户信息，有 Redis，user + cnt
	// 注意 is_follow 字段默认值，需要其他服务调用 favorite 服务取得
	UserInfo(context.Context, *UserRequest) (*UserResponse, error)
	// UserInfoList 获取用户信息列表，使用用户 ID 列表查询
	UserInfoList(context.Context, *UserInfoListrRequest) (*UserInfoListResponse, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/douyin/user/register", _UserService_Register0_HTTP_Handler(srv))
	r.POST("/douyin/user/login", _UserService_Login0_HTTP_Handler(srv))
	r.GET("/douyin/user", _UserService_UserInfo0_HTTP_Handler(srv))
	r.POST("/douyin/user/update_work_count", _UserService_UpdateWorkCnt0_HTTP_Handler(srv))
	r.POST("/douyin/user/update_favorite_count", _UserService_UpdateFavoriteCnt0_HTTP_Handler(srv))
	r.POST("/douyin/user/update_follow_count", _UserService_UpdateFollowCnt0_HTTP_Handler(srv))
	r.POST("/douyin/user/list", _UserService_UserInfoList0_HTTP_Handler(srv))
}

func _UserService_Register0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserRegisterRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceRegister)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Register(ctx, req.(*UserRegisterRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserRegisterResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_Login0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserLoginRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*UserLoginRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserLoginResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_UserInfo0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUserInfo)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UserInfo(ctx, req.(*UserRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateWorkCnt0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateWorkCntRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateWorkCnt)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateWorkCnt(ctx, req.(*UpdateWorkCntRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateWorkCntResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateFavoriteCnt0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateFavoriteCntRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateFavoriteCnt)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateFavoriteCnt(ctx, req.(*UpdateFavoriteCntRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateFavoriteCntResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateFollowCnt0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateFollowCntRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateFollowCnt)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateFollowCnt(ctx, req.(*UpdateFollowCntRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UpdateFollowCntResponse)
		return ctx.Result(200, reply)
	}
}

func _UserService_UserInfoList0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserInfoListrRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUserInfoList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UserInfoList(ctx, req.(*UserInfoListrRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserInfoListResponse)
		return ctx.Result(200, reply)
	}
}

type UserServiceHTTPClient interface {
	Login(ctx context.Context, req *UserLoginRequest, opts ...http.CallOption) (rsp *UserLoginResponse, err error)
	Register(ctx context.Context, req *UserRegisterRequest, opts ...http.CallOption) (rsp *UserRegisterResponse, err error)
	UpdateFavoriteCnt(ctx context.Context, req *UpdateFavoriteCntRequest, opts ...http.CallOption) (rsp *UpdateFavoriteCntResponse, err error)
	UpdateFollowCnt(ctx context.Context, req *UpdateFollowCntRequest, opts ...http.CallOption) (rsp *UpdateFollowCntResponse, err error)
	UpdateWorkCnt(ctx context.Context, req *UpdateWorkCntRequest, opts ...http.CallOption) (rsp *UpdateWorkCntResponse, err error)
	UserInfo(ctx context.Context, req *UserRequest, opts ...http.CallOption) (rsp *UserResponse, err error)
	UserInfoList(ctx context.Context, req *UserInfoListrRequest, opts ...http.CallOption) (rsp *UserInfoListResponse, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) Login(ctx context.Context, in *UserLoginRequest, opts ...http.CallOption) (*UserLoginResponse, error) {
	var out UserLoginResponse
	pattern := "/douyin/user/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) Register(ctx context.Context, in *UserRegisterRequest, opts ...http.CallOption) (*UserRegisterResponse, error) {
	var out UserRegisterResponse
	pattern := "/douyin/user/register"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceRegister))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateFavoriteCnt(ctx context.Context, in *UpdateFavoriteCntRequest, opts ...http.CallOption) (*UpdateFavoriteCntResponse, error) {
	var out UpdateFavoriteCntResponse
	pattern := "/douyin/user/update_favorite_count"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateFavoriteCnt))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateFollowCnt(ctx context.Context, in *UpdateFollowCntRequest, opts ...http.CallOption) (*UpdateFollowCntResponse, error) {
	var out UpdateFollowCntResponse
	pattern := "/douyin/user/update_follow_count"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateFollowCnt))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateWorkCnt(ctx context.Context, in *UpdateWorkCntRequest, opts ...http.CallOption) (*UpdateWorkCntResponse, error) {
	var out UpdateWorkCntResponse
	pattern := "/douyin/user/update_work_count"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateWorkCnt))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UserInfo(ctx context.Context, in *UserRequest, opts ...http.CallOption) (*UserResponse, error) {
	var out UserResponse
	pattern := "/douyin/user"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationUserServiceUserInfo))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UserInfoList(ctx context.Context, in *UserInfoListrRequest, opts ...http.CallOption) (*UserInfoListResponse, error) {
	var out UserInfoListResponse
	pattern := "/douyin/user/list"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUserInfoList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
