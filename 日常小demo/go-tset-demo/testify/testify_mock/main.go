//参考教程：https://mp.weixin.qq.com/s/4H66WsEHaDt5enkROiYgOQ
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
  Name string
  Age  int
}

//接口，提供获取用户列表的方法
type ICrawler interface {
  GetUserList() ([]*User, error)
}

//实现GetUserList接口
type MyCrawler struct {
  url string
}

func (c *MyCrawler) GetUserList() ([]*User, error) {
  resp, err := http.Get(c.url)
  if err != nil {
    return nil, err
  }

  defer resp.Body.Close()
  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }

  var userList []*User
  err = json.Unmarshal(data, &userList)
  if err != nil {
    return nil, err
  }

  return userList, nil
}

//需要测试的函数，传入的是接口类型，才可以进行mock
func GetAndPrintUsers(crawler ICrawler) {
  users, err := crawler.GetUserList()
  if err != nil {
    return
  }

  for _, u := range users {
    fmt.Println(u)
  }
}