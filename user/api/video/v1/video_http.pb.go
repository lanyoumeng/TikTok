// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v3.12.4
// source: video/v1/video.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	"io/ioutil"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationVideoServiceFavoriteListByVId = "/video.api.video.v1.VideoService/FavoriteListByVId"
const OperationVideoServiceFeed = "/video.api.video.v1.VideoService/Feed"
const OperationVideoServiceGetAIdByVId = "/video.api.video.v1.VideoService/GetAIdByVId"
const OperationVideoServicePublish = "/video.api.video.v1.VideoService/Publish"
const OperationVideoServicePublishList = "/video.api.video.v1.VideoService/PublishList"
const OperationVideoServicePublishVidsByAId = "/video.api.video.v1.VideoService/PublishVidsByAId"
const OperationVideoServiceWorkCnt = "/video.api.video.v1.VideoService/WorkCnt"

type VideoServiceHTTPServer interface {
	// FavoriteListByVId 通过喜欢视频 id 列表获取视频列表
	FavoriteListByVId(context.Context, *FavoriteListReq) (*FavoriteListResp, error)
	// Feed 获取视频流
	Feed(context.Context, *DouyinFeedRequest) (*DouyinFeedResponse, error)
	// GetAIdByVId 通过视频 id 获取作者 id
	GetAIdByVId(context.Context, *GetAIdByVIdReq) (*GetAIdByVIdResp, error)
	// Publish 发布视频
	Publish(context.Context, *DouyinPublishActionRequest) (*DouyinPublishActionResponse, error)
	// PublishList 获取用户发布的视频列表
	PublishList(context.Context, *DouyinPublishListRequest) (*DouyinPublishListResponse, error)
	// PublishVidsByAId 通过作者 id 获取作者的发布视频 id 列表
	PublishVidsByAId(context.Context, *PublishVidsByAIdReq) (*PublishVidsByAIdResp, error)
	// WorkCnt 获取用户作品数量
	WorkCnt(context.Context, *WorkCntRequest) (*WorkCntResponse, error)
}

func RegisterVideoServiceHTTPServer(s *http.Server, srv VideoServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/douyin/feed", _VideoService_Feed0_HTTP_Handler(srv))
	r.POST("/douyin/publish/action", _VideoService_Publish0_HTTP_Handler(srv))
	r.GET("/douyin/publish/list", _VideoService_PublishList0_HTTP_Handler(srv))
	r.GET("/douyin/work/count", _VideoService_WorkCnt0_HTTP_Handler(srv))
	r.POST("/douyin/favorite/list", _VideoService_FavoriteListByVId0_HTTP_Handler(srv))
	r.GET("/douyin/publish/vids", _VideoService_PublishVidsByAId0_HTTP_Handler(srv))
	r.GET("/douyin/author/video", _VideoService_GetAIdByVId0_HTTP_Handler(srv))
}

func _VideoService_Feed0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinFeedRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServiceFeed)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Feed(ctx, req.(*DouyinFeedRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinFeedResponse)
		return ctx.Result(200, reply)
	}
}

func _VideoService_Publish0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		// 解析 multipart/form-data 请求
		if err := ctx.Request().ParseMultipartForm(32 << 20); err != nil {
			return err
		}

		// 处理文件字段
		file, _, err := ctx.Request().FormFile("data")
		if err != nil {
			return err
		}
		defer file.Close()

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		// 获取其他字段
		token := ctx.Request().FormValue("token")
		title := ctx.Request().FormValue("title")

		// 创建 gRPC 请求对象
		in := &DouyinPublishActionRequest{
			Data: fileBytes,
			Token:    token,
			Title:    title,
		}

		http.SetOperation(ctx, OperationVideoServicePublish)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Publish(ctx, req.(*DouyinPublishActionRequest))
		})
		out, err := h(ctx, in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinPublishActionResponse)
		return ctx.Result(200, reply)
	}
}

func _VideoService_PublishList0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinPublishListRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServicePublishList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PublishList(ctx, req.(*DouyinPublishListRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinPublishListResponse)
		return ctx.Result(200, reply)
	}
}

func _VideoService_WorkCnt0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in WorkCntRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServiceWorkCnt)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.WorkCnt(ctx, req.(*WorkCntRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*WorkCntResponse)
		return ctx.Result(200, reply)
	}
}

func _VideoService_FavoriteListByVId0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in FavoriteListReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServiceFavoriteListByVId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.FavoriteListByVId(ctx, req.(*FavoriteListReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*FavoriteListResp)
		return ctx.Result(200, reply)
	}
}

func _VideoService_PublishVidsByAId0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PublishVidsByAIdReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServicePublishVidsByAId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.PublishVidsByAId(ctx, req.(*PublishVidsByAIdReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PublishVidsByAIdResp)
		return ctx.Result(200, reply)
	}
}

func _VideoService_GetAIdByVId0_HTTP_Handler(srv VideoServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetAIdByVIdReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationVideoServiceGetAIdByVId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetAIdByVId(ctx, req.(*GetAIdByVIdReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetAIdByVIdResp)
		return ctx.Result(200, reply)
	}
}

type VideoServiceHTTPClient interface {
	FavoriteListByVId(ctx context.Context, req *FavoriteListReq, opts ...http.CallOption) (rsp *FavoriteListResp, err error)
	Feed(ctx context.Context, req *DouyinFeedRequest, opts ...http.CallOption) (rsp *DouyinFeedResponse, err error)
	GetAIdByVId(ctx context.Context, req *GetAIdByVIdReq, opts ...http.CallOption) (rsp *GetAIdByVIdResp, err error)
	Publish(ctx context.Context, req *DouyinPublishActionRequest, opts ...http.CallOption) (rsp *DouyinPublishActionResponse, err error)
	PublishList(ctx context.Context, req *DouyinPublishListRequest, opts ...http.CallOption) (rsp *DouyinPublishListResponse, err error)
	PublishVidsByAId(ctx context.Context, req *PublishVidsByAIdReq, opts ...http.CallOption) (rsp *PublishVidsByAIdResp, err error)
	WorkCnt(ctx context.Context, req *WorkCntRequest, opts ...http.CallOption) (rsp *WorkCntResponse, err error)
}

type VideoServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewVideoServiceHTTPClient(client *http.Client) VideoServiceHTTPClient {
	return &VideoServiceHTTPClientImpl{client}
}

func (c *VideoServiceHTTPClientImpl) FavoriteListByVId(ctx context.Context, in *FavoriteListReq, opts ...http.CallOption) (*FavoriteListResp, error) {
	var out FavoriteListResp
	pattern := "/douyin/favorite/list"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationVideoServiceFavoriteListByVId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) Feed(ctx context.Context, in *DouyinFeedRequest, opts ...http.CallOption) (*DouyinFeedResponse, error) {
	var out DouyinFeedResponse
	pattern := "/douyin/feed"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationVideoServiceFeed))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) GetAIdByVId(ctx context.Context, in *GetAIdByVIdReq, opts ...http.CallOption) (*GetAIdByVIdResp, error) {
	var out GetAIdByVIdResp
	pattern := "/douyin/author/video"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationVideoServiceGetAIdByVId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) Publish(ctx context.Context, in *DouyinPublishActionRequest, opts ...http.CallOption) (*DouyinPublishActionResponse, error) {
	var out DouyinPublishActionResponse
	pattern := "/douyin/publish/action"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationVideoServicePublish))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) PublishList(ctx context.Context, in *DouyinPublishListRequest, opts ...http.CallOption) (*DouyinPublishListResponse, error) {
	var out DouyinPublishListResponse
	pattern := "/douyin/publish/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationVideoServicePublishList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) PublishVidsByAId(ctx context.Context, in *PublishVidsByAIdReq, opts ...http.CallOption) (*PublishVidsByAIdResp, error) {
	var out PublishVidsByAIdResp
	pattern := "/douyin/publish/vids"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationVideoServicePublishVidsByAId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *VideoServiceHTTPClientImpl) WorkCnt(ctx context.Context, in *WorkCntRequest, opts ...http.CallOption) (*WorkCntResponse, error) {
	var out WorkCntResponse
	pattern := "/douyin/work/count"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationVideoServiceWorkCnt))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
