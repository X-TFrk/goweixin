package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"syscall"
	"testing"
)

func TestMiniProgramClient_SendMessage(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	appId, _ = syscall.Getenv("appId")
	appSecret, _ = syscall.Getenv("appSecret")

	fmt.Println(appId, appSecret)

	c := NewMiniProgramClient(appId, appSecret, "")

	token, err := c.AuthGetAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("token is:", token)

	openId := "on9VO5YXH_gMLxRKMql98IUjtzkI"
	templateId := "IgOxNz7ydQn9UiU09IgJtENIpIwigg5TroAbRLXcosY"
	page := "index?foo=bar"
	data := map[string]string{"thing1": "这是一个内容", "thing7": "这个也是内容", "thing3": "这个也是内容啊"}

	err = c.SendMessage(token, openId, templateId, page, data)
	if err != nil {
		fmt.Println("send err:", err.Error())
		return
	}

	err = c.SendMessage(token, openId, templateId, page, data)
	if err != nil {
		fmt.Println("send err:", err.Error())
		return
	}
}
