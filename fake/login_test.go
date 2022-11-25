package fake

import (
	"fmt"
	"testing"
)

func TestLogin(t *testing.T) {
	appId := ""
	appSecret := ""
	code := "xxx"
	info, err := Login(appId, appSecret, code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(info)
}
