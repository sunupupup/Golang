package main

import (
	"testing"

	"github.com/golang/mock/gomock"
)

//go test . -cover -v
func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish() //断言 DB.Get() 方法是否被调用，不然下面的代码没有意义

	m := NewMockDB(ctrl)
	//m.EXPECT().Get("sjw").Return(25, errors.New("not exist"))
	m.EXPECT().Get("sjw").Return(25, nil)

	//如果 DB.Get() 返回 error，那么 GetFromDB() 返回 -1)。
	if v := GetFromDB(m, "sjw"); v != 25 {
		t.Fatal("expected -1, but got ", v)
	}

	m.EXPECT().Get(gomock.Eq("tao")).Return(18, nil)
	if v := GetFromDB(m, "tao"); v != 18 {
		t.Fatal("expected 18, but got ", v)

	}
}
