package goweixin

import (
	"fmt"
	"testing"
)

func TestMiniLogin(t *testing.T) {
	appId := ""
	appSecret := ""
	code := "xxx"
	encryptedData := "加密的数据"
	iv := "ssss"
	userInfo, err := MiniLogin(appId, appSecret, code, encryptedData, iv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(userInfo)
}
