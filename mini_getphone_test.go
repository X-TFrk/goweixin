package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"syscall"
	"testing"
)

func TestMiniProgramClient_GetPhoneNumber(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	appId, _ = syscall.Getenv("appId")
	appSecret, _ = syscall.Getenv("appSecret")

	fmt.Println(appId, appSecret)

	c := NewMiniProgramClient(appId, appSecret)

	token, err := c.AuthGetAccessToken()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("token is:", token)

	code := "910031e46a34e633401c2ebb23f281646ea9775ad8c1276b793e59846f0ddb22"
	phone, err := c.GetPhoneNumber(token, code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%#v", phone)
}
