// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.3
// - protoc             v3.12.4
// source: comment/v1/comment.proto

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

const OperationCommentServiceComment = "/comment.api.comment.v1.CommentService/Comment"
const OperationCommentServiceCommentList = "/comment.api.comment.v1.CommentService/CommentList"
const OperationCommentServiceGetCommentCntByVId = "/comment.api.comment.v1.CommentService/GetCommentCntByVId"

type CommentServiceHTTPServer interface {
	// Comment 发布/删除评论
	Comment(context.Context, *DouyinCommentSendRequest) (*DouyinCommentSendResponse, error)
	// CommentList 获取评论列表
	CommentList(context.Context, *DouyinCommentListRequest) (*DouyinCommentListResponse, error)
	// GetCommentCntByVId 获取视频评论总数
	GetCommentCntByVId(context.Context, *GetCommentCntByVIdReq) (*GetCommentCntByVIdResp, error)
}

func RegisterCommentServiceHTTPServer(s *http.Server, srv CommentServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/douyin/comment/action", _CommentService_Comment0_HTTP_Handler(srv))
	r.GET("/douyin/comment/list", _CommentService_CommentList0_HTTP_Handler(srv))
	r.GET("/douyin/comment/count/video", _CommentService_GetCommentCntByVId0_HTTP_Handler(srv))
}

func _CommentService_Comment0_HTTP_Handler(srv CommentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinCommentSendRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCommentServiceComment)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Comment(ctx, req.(*DouyinCommentSendRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinCommentSendResponse)
		return ctx.Result(200, reply)
	}
}

func _CommentService_CommentList0_HTTP_Handler(srv CommentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinCommentListRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCommentServiceCommentList)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.CommentList(ctx, req.(*DouyinCommentListRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinCommentListResponse)
		return ctx.Result(200, reply)
	}
}

func _CommentService_GetCommentCntByVId0_HTTP_Handler(srv CommentServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetCommentCntByVIdReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationCommentServiceGetCommentCntByVId)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCommentCntByVId(ctx, req.(*GetCommentCntByVIdReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetCommentCntByVIdResp)
		return ctx.Result(200, reply)
	}
}

type CommentServiceHTTPClient interface {
	Comment(ctx context.Context, req *DouyinCommentSendRequest, opts ...http.CallOption) (rsp *DouyinCommentSendResponse, err error)
	CommentList(ctx context.Context, req *DouyinCommentListRequest, opts ...http.CallOption) (rsp *DouyinCommentListResponse, err error)
	GetCommentCntByVId(ctx context.Context, req *GetCommentCntByVIdReq, opts ...http.CallOption) (rsp *GetCommentCntByVIdResp, err error)
}

type CommentServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewCommentServiceHTTPClient(client *http.Client) CommentServiceHTTPClient {
	return &CommentServiceHTTPClientImpl{client}
}

func (c *CommentServiceHTTPClientImpl) Comment(ctx context.Context, in *DouyinCommentSendRequest, opts ...http.CallOption) (*DouyinCommentSendResponse, error) {
	var out DouyinCommentSendResponse
	pattern := "/douyin/comment/action"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationCommentServiceComment))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CommentServiceHTTPClientImpl) CommentList(ctx context.Context, in *DouyinCommentListRequest, opts ...http.CallOption) (*DouyinCommentListResponse, error) {
	var out DouyinCommentListResponse
	pattern := "/douyin/comment/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCommentServiceCommentList))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *CommentServiceHTTPClientImpl) GetCommentCntByVId(ctx context.Context, in *GetCommentCntByVIdReq, opts ...http.CallOption) (*GetCommentCntByVIdResp, error) {
	var out GetCommentCntByVIdResp
	pattern := "/douyin/comment/count/video"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationCommentServiceGetCommentCntByVId))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
