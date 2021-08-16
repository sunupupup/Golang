package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

//对接口进行mock
type MockCrawler struct {
	mock.Mock
}

//实现接口的GetUserList方法
func (m *MockCrawler) GetUserList() ([]*User, error) {
	args := m.Called()
	return args.Get(0).([]*User), args.Error(1)
}

//
var (
	MockUsers []*User
)

func init() {
	MockUsers = append(MockUsers, &User{"dj", 18})
	MockUsers = append(MockUsers, &User{"zhangsan", 20})
}

func TestGetUserList(t *testing.T) {
	//创建mock对象
	crawler := new(MockCrawler)
	//这里指示调用GetUserList()方法的返回值分别为MockUsers和nil，就是模拟返回结果
	//返回的值被return args.Get(0).([]*User), args.Error(1)这两个返回值捕获
	crawler.On("GetUserList").Return(MockUsers, nil)

	//测试GetAndPrintUsers函数
	GetAndPrintUsers(crawler)

	crawler.AssertExpectations(t)
}
