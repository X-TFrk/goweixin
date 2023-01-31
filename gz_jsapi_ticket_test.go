package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"syscall"
	"testing"
)

func TestGzClient_GetJsapiTicketAndSign(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	appId, _ = syscall.Getenv("appId")
	appSecret, _ = syscall.Getenv("appSecret")

	fmt.Println(appId, appSecret)

	c := NewGzClient(appId, appSecret)

	result, sign, err := c.GetJsapiTicketAndSign("http://mp.weixin.qq.com?params=value")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%#v\n%#v\n", result, sign)
	fmt.Printf("%#v", c)
}
