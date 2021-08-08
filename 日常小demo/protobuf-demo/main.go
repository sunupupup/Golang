package main

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

//protoc --go_out=. addressbook.proto  生成对应的go文件
//go mod tidy

//go run main.go addressbook.pb.go  需要将addressbook.pb.go加上才可通过编译
func main() {
	me := &Person{
		Name: "sjw",
		Age:  25,
		SocialFollowers: &SocialFollowers{
			YouTube: 2500,
			Twitter: 1400,
		},
	}

	data, err := proto.Marshal(me) //protobuf编码
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	you := &Person{}
	err = proto.Unmarshal(data, you) //protobuf解码
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(you.Name, you.Age, you.GetSocialFollowers().Twitter, you.GetSocialFollowers().GetYouTube())
}
