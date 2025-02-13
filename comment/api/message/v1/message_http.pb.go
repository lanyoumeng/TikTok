// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.3
// - protoc             v3.12.4
// source: message/v1/message.proto

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

const OperationMessageServiceGetNewMessages = "/message.api.message.v1.MessageService/GetNewMessages"
const OperationMessageServiceMessageRecord = "/message.api.message.v1.MessageService/MessageRecord"
const OperationMessageServiceMessageSend = "/message.api.message.v1.MessageService/MessageSend"

type MessageServiceHTTPServer interface {
	// GetNewMessages 获取用户和每个好友的最新一条消息
	GetNewMessages(context.Context, *GetNewMessagesRequest) (*GetNewMessagesResponse, error)
	// MessageRecord 获取历史消息记录
	MessageRecord(context.Context, *DouyinMessageRecordRequest) (*DouyinMessageRecordResponse, error)
	// MessageSend 发送消息
	MessageSend(context.Context, *DouyinMessageSendRequest) (*DouyinMessageSendResponse, error)
}

func RegisterMessageServiceHTTPServer(s *http.Server, srv MessageServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/douyin/message/chat", _MessageService_MessageRecord0_HTTP_Handler(srv))
	r.POST("/douyin/message/action", _MessageService_MessageSend0_HTTP_Handler(srv))
	r.GET("/douyin/message/latest", _MessageService_GetNewMessages0_HTTP_Handler(srv))
}

func _MessageService_MessageRecord0_HTTP_Handler(srv MessageServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinMessageRecordRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMessageServiceMessageRecord)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.MessageRecord(ctx, req.(*DouyinMessageRecordRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinMessageRecordResponse)
		return ctx.Result(200, reply)
	}
}

func _MessageService_MessageSend0_HTTP_Handler(srv MessageServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DouyinMessageSendRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMessageServiceMessageSend)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.MessageSend(ctx, req.(*DouyinMessageSendRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*DouyinMessageSendResponse)
		return ctx.Result(200, reply)
	}
}

func _MessageService_GetNewMessages0_HTTP_Handler(srv MessageServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetNewMessagesRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationMessageServiceGetNewMessages)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetNewMessages(ctx, req.(*GetNewMessagesRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetNewMessagesResponse)
		return ctx.Result(200, reply)
	}
}

type MessageServiceHTTPClient interface {
	GetNewMessages(ctx context.Context, req *GetNewMessagesRequest, opts ...http.CallOption) (rsp *GetNewMessagesResponse, err error)
	MessageRecord(ctx context.Context, req *DouyinMessageRecordRequest, opts ...http.CallOption) (rsp *DouyinMessageRecordResponse, err error)
	MessageSend(ctx context.Context, req *DouyinMessageSendRequest, opts ...http.CallOption) (rsp *DouyinMessageSendResponse, err error)
}

type MessageServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewMessageServiceHTTPClient(client *http.Client) MessageServiceHTTPClient {
	return &MessageServiceHTTPClientImpl{client}
}

func (c *MessageServiceHTTPClientImpl) GetNewMessages(ctx context.Context, in *GetNewMessagesRequest, opts ...http.CallOption) (*GetNewMessagesResponse, error) {
	var out GetNewMessagesResponse
	pattern := "/douyin/message/latest"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationMessageServiceGetNewMessages))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *MessageServiceHTTPClientImpl) MessageRecord(ctx context.Context, in *DouyinMessageRecordRequest, opts ...http.CallOption) (*DouyinMessageRecordResponse, error) {
	var out DouyinMessageRecordResponse
	pattern := "/douyin/message/chat"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationMessageServiceMessageRecord))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *MessageServiceHTTPClientImpl) MessageSend(ctx context.Context, in *DouyinMessageSendRequest, opts ...http.CallOption) (*DouyinMessageSendResponse, error) {
	var out DouyinMessageSendResponse
	pattern := "/douyin/message/action"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationMessageServiceMessageSend))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
