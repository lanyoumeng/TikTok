// Code generated by MockGen. DO NOT EDIT.
// Source: video/internal/biz (interfaces: VideoRepo)

// Package mrepo is a generated GoMock package.
package mrepo

import (
	context "context"
	reflect "reflect"
	time "time"
	v1 "video/api/video/v1"
	model "video/internal/pkg/model"

	gomock "github.com/golang/mock/gomock"
)

// MockVideoRepo is a mock of VideoRepo interface.
type MockVideoRepo struct {
	ctrl     *gomock.Controller
	recorder *MockVideoRepoMockRecorder
}

// MockVideoRepoMockRecorder is the mock recorder for MockVideoRepo.
type MockVideoRepoMockRecorder struct {
	mock *MockVideoRepo
}

// NewMockVideoRepo creates a new mock instance.
func NewMockVideoRepo(ctrl *gomock.Controller) *MockVideoRepo {
	mock := &MockVideoRepo{ctrl: ctrl}
	mock.recorder = &MockVideoRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVideoRepo) EXPECT() *MockVideoRepoMockRecorder {
	return m.recorder
}

// GetAuthorInfoById mocks base method.
func (m *MockVideoRepo) GetAuthorInfoById(arg0 context.Context, arg1 int64) (*v1.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAuthorInfoById", arg0, arg1)
	ret0, _ := ret[0].(*v1.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAuthorInfoById indicates an expected call of GetAuthorInfoById.
func (mr *MockVideoRepoMockRecorder) GetAuthorInfoById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAuthorInfoById", reflect.TypeOf((*MockVideoRepo)(nil).GetAuthorInfoById), arg0, arg1)
}

// GetCommentCntByVId mocks base method.
func (m *MockVideoRepo) GetCommentCntByVId(arg0 context.Context, arg1 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommentCntByVId", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommentCntByVId indicates an expected call of GetCommentCntByVId.
func (mr *MockVideoRepoMockRecorder) GetCommentCntByVId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommentCntByVId", reflect.TypeOf((*MockVideoRepo)(nil).GetCommentCntByVId), arg0, arg1)
}

// GetFavoriteCntByVId mocks base method.
func (m *MockVideoRepo) GetFavoriteCntByVId(arg0 context.Context, arg1 int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavoriteCntByVId", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoriteCntByVId indicates an expected call of GetFavoriteCntByVId.
func (mr *MockVideoRepoMockRecorder) GetFavoriteCntByVId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoriteCntByVId", reflect.TypeOf((*MockVideoRepo)(nil).GetFavoriteCntByVId), arg0, arg1)
}

// GetIsFavorite mocks base method.
func (m *MockVideoRepo) GetIsFavorite(arg0 context.Context, arg1, arg2 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIsFavorite", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIsFavorite indicates an expected call of GetIsFavorite.
func (mr *MockVideoRepoMockRecorder) GetIsFavorite(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIsFavorite", reflect.TypeOf((*MockVideoRepo)(nil).GetIsFavorite), arg0, arg1, arg2)
}

// GetVideoListByAuthorId mocks base method.
func (m *MockVideoRepo) GetVideoListByAuthorId(arg0 context.Context, arg1 int64) ([]*model.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoListByAuthorId", arg0, arg1)
	ret0, _ := ret[0].([]*model.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoListByAuthorId indicates an expected call of GetVideoListByAuthorId.
func (mr *MockVideoRepoMockRecorder) GetVideoListByAuthorId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoListByAuthorId", reflect.TypeOf((*MockVideoRepo)(nil).GetVideoListByAuthorId), arg0, arg1)
}

// GetVideoListByLatestTime mocks base method.
func (m *MockVideoRepo) GetVideoListByLatestTime(arg0 context.Context, arg1 time.Time) ([]*model.Video, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoListByLatestTime", arg0, arg1)
	ret0, _ := ret[0].([]*model.Video)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetVideoListByLatestTime indicates an expected call of GetVideoListByLatestTime.
func (mr *MockVideoRepoMockRecorder) GetVideoListByLatestTime(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoListByLatestTime", reflect.TypeOf((*MockVideoRepo)(nil).GetVideoListByLatestTime), arg0, arg1)
}

// GetvideoByVId mocks base method.
func (m *MockVideoRepo) GetvideoByVId(arg0 context.Context, arg1 int64) (*model.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetvideoByVId", arg0, arg1)
	ret0, _ := ret[0].(*model.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetvideoByVId indicates an expected call of GetvideoByVId.
func (mr *MockVideoRepoMockRecorder) GetvideoByVId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetvideoByVId", reflect.TypeOf((*MockVideoRepo)(nil).GetvideoByVId), arg0, arg1)
}

// PublishKafka mocks base method.
func (m *MockVideoRepo) PublishKafka(arg0 context.Context, arg1 *model.VideoKafkaMessage) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PublishKafka", arg0, arg1)
}

// PublishKafka indicates an expected call of PublishKafka.
func (mr *MockVideoRepoMockRecorder) PublishKafka(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishKafka", reflect.TypeOf((*MockVideoRepo)(nil).PublishKafka), arg0, arg1)
}

// RGetVideoInfo mocks base method.
func (m *MockVideoRepo) RGetVideoInfo(arg0 context.Context, arg1 int64) (*model.Video, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetVideoInfo", arg0, arg1)
	ret0, _ := ret[0].(*model.Video)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetVideoInfo indicates an expected call of RGetVideoInfo.
func (mr *MockVideoRepoMockRecorder) RGetVideoInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetVideoInfo", reflect.TypeOf((*MockVideoRepo)(nil).RGetVideoInfo), arg0, arg1)
}

// RPublishVidsByAuthorId mocks base method.
func (m *MockVideoRepo) RPublishVidsByAuthorId(arg0 context.Context, arg1 int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RPublishVidsByAuthorId", arg0, arg1)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RPublishVidsByAuthorId indicates an expected call of RPublishVidsByAuthorId.
func (mr *MockVideoRepoMockRecorder) RPublishVidsByAuthorId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RPublishVidsByAuthorId", reflect.TypeOf((*MockVideoRepo)(nil).RPublishVidsByAuthorId), arg0, arg1)
}

// RSavePublishVids mocks base method.
func (m *MockVideoRepo) RSavePublishVids(arg0 context.Context, arg1 int64, arg2 []int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RSavePublishVids", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RSavePublishVids indicates an expected call of RSavePublishVids.
func (mr *MockVideoRepoMockRecorder) RSavePublishVids(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RSavePublishVids", reflect.TypeOf((*MockVideoRepo)(nil).RSavePublishVids), arg0, arg1, arg2)
}

// RSaveVideoInfo mocks base method.
func (m *MockVideoRepo) RSaveVideoInfo(arg0 context.Context, arg1 *model.Video) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RSaveVideoInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RSaveVideoInfo indicates an expected call of RSaveVideoInfo.
func (mr *MockVideoRepoMockRecorder) RSaveVideoInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RSaveVideoInfo", reflect.TypeOf((*MockVideoRepo)(nil).RSaveVideoInfo), arg0, arg1)
}

// RZSetSaveVIds mocks base method.
func (m *MockVideoRepo) RZSetSaveVIds(arg0 context.Context, arg1 []int64, arg2 []float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RZSetSaveVIds", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RZSetSaveVIds indicates an expected call of RZSetSaveVIds.
func (mr *MockVideoRepoMockRecorder) RZSetSaveVIds(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RZSetSaveVIds", reflect.TypeOf((*MockVideoRepo)(nil).RZSetSaveVIds), arg0, arg1, arg2)
}

// RZSetVideoIds mocks base method.
func (m *MockVideoRepo) RZSetVideoIds(arg0 context.Context, arg1 string, arg2 int64) (int64, []int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RZSetVideoIds", arg0, arg1, arg2)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].([]int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// RZSetVideoIds indicates an expected call of RZSetVideoIds.
func (mr *MockVideoRepoMockRecorder) RZSetVideoIds(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RZSetVideoIds", reflect.TypeOf((*MockVideoRepo)(nil).RZSetVideoIds), arg0, arg1, arg2)
}

// Save mocks base method.
func (m *MockVideoRepo) Save(arg0 context.Context, arg1 *model.Video) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockVideoRepoMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockVideoRepo)(nil).Save), arg0, arg1)
}