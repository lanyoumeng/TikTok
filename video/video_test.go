package mrepo_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "video/api/video/v1"
	"video/internal/mocks/mrepo"
)

func TestGetAuthorInfoById(t *testing.T) {
	// 创建一个新的控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 创建一个新的 MockVideoRepo 实例
	mockRepo := mrepo.NewMockVideoRepo(ctrl)

	// 设置期望值和返回值
	expectedUser := &v1.User{Id: 1, Name: "Author"}
	mockRepo.EXPECT().GetAuthorInfoById(gomock.Any(), int64(1)).Return(expectedUser, nil)

	// 调用被测试的函数
	user, err := mockRepo.GetAuthorInfoById(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Id != expectedUser.Id || user.Name != expectedUser.Name {
		t.Fatalf("expected %v, got %v", expectedUser, user)
	}
}
