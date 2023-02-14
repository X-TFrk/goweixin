package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"testing"
)

func TestOpenClient_LoginGetBaseInfo(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxbdc5610cc59c1631"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	fmt.Println(appId, appSecret)

	c := NewOpenClient(appId, appSecret)

	code := "0516D5Ga1pW3ME07bzHa1zS8Jk36D5GK"
	result, err := c.LoginGetBaseInfo(code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s,%s", result.OpenId, result.UnionId)
}

func TestOpenClient_LoginGetUserInfo(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxbdc5610cc59c1631"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	fmt.Println(appId, appSecret)

	c := NewOpenClient(appId, appSecret)

	code := "071yEk100YW0tP15wg000TLA6r2yEk1S"
	result, err := c.LoginGetUserInfo(code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s,%s\n", result.OpenId, result.UnionId)
	fmt.Printf("%s,%s", result.NickName, result.Img)
}
