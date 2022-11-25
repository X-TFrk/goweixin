package goweixin

import (
	"fmt"
	"testing"
)

func TestMiniProgramClient_SendMessage(t *testing.T) {
	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	c := NewMiniProgramClient(appId, appSecret)

	token, err := c.AuthGetAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("token is:", token)

	openId := "omvdI6yayVSRLK9NL2OcCHEWQ0mA"
	templateId := "templateId"
	page := ""
	data := map[string]string{"thing1": "2222", "thing7": "sss", "thing3": "dddd"}
	state := MiniProgramStateDeveloper

	err = c.SendMessage(token, openId, templateId, page, data, state)
	if err != nil {
		fmt.Println("send err:", err.Error())
		return
	}
}
