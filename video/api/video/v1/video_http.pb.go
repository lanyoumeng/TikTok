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

const OperationVideoServicePublish = "/video.api.video.v1.VideoService/Publish"

type VideoServiceHTTPServer interface {
	Publish(context.Context, *DouyinPublishActionRequest) (*DouyinPublishActionResponse, error)
}

func RegisterVideoServiceHTTPServer(s *http.Server, srv VideoServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/douyin/publish/action", _VideoService_Publish0_HTTP_Handler(srv))
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

type VideoServiceHTTPClient interface {
	Publish(ctx context.Context, req *DouyinPublishActionRequest, opts ...http.CallOption) (rsp *DouyinPublishActionResponse, err error)
}

type VideoServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewVideoServiceHTTPClient(client *http.Client) VideoServiceHTTPClient {
	return &VideoServiceHTTPClientImpl{client}
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
