package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"testing"
)

func TestSendMessage(t *testing.T) {
	miner.SetLogLevel(miner.DEBUG)
	appId := ""
	appSecret := ""
	token, err := AuthGetAccessToken(appId, appSecret)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("token is:", token)

	openId := "sss"
	templateId := ""
	page := ""
	data := map[string]string{"thing1": "2222", "thing7": "sss", "thing3": "dddd"}
	state := wxStateFormal

	err = SendMessage(token, openId, templateId, page, data, state)
	if err != nil {
		fmt.Println("send err:", err.Error())
		return
	}
}
