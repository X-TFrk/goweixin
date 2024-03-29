package goweixin

import (
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"syscall"
	"testing"
)

func TestMiniProgramClient_LoginGetBaseInfo(t *testing.T) {
	miner.SetLogLevel(miner.INFO)

	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	appId, _ = syscall.Getenv("appId")
	appSecret, _ = syscall.Getenv("appSecret")

	fmt.Println(appId, appSecret)

	c := NewMiniProgramClient(appId, appSecret, "")

	code := "033SD2100a5hZO15kd100qqPJf2SD21f"
	result, err := c.LoginGetBaseInfo(code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%s,%s", result.OpenId, result.UnionId)
}

func TestMiniProgramClient_LoginGetUserInfo(t *testing.T) {
	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	c := NewMiniProgramClient(appId, appSecret, "")

	code := "013TX31w326TEZ2I2f4w3Hol9j2TX31N"
	encryptedData := "98A19r2TH/F+biCYbx6YE9dnMjWVZfUgqPEWFenYX3jP8JpIKihNCyjE/Or0/pmYT+PCn6wCmV7s5LDwwQ92kcrMpuInOrmPWD36pI0mfywqh+53BcN4G+d30aG6ehCV3hPEqxE35ImpXxE5xuWqqsLX0YvgCgA5hLBGRWRPGiiXN4eSrLvCNI58BlC8VG16Iz6Z89NqAQ5WsQrWPjJqygZOsGnkvTFTwKzs6eM4jlmceWv4B37NxecCAwRkkHZRxla5mdluL5lwKovfH5feDQTg2Ui2/4Mc/raIh9tXV6lUqEkn8f4yJFhjXwJuhXdITxLcKnqP6O/n1DUOn1xNh5imZEiDBil14Zvy71pvb5JdRPhurtqOFM/5wCRc7BXwrrDo1iMUw3e7OOL9gFu7Fg=="
	iv := "+DTYOLIe/s1ObkVKZQolpw=="

	userInfo, err := c.LoginGetUserInfo(code, encryptedData, iv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%#v", userInfo)
}

func TestMiniProgramClient_DecryptUserInfo(t *testing.T) {
	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"
	sessionKey := "xmOgG3F4QUh4uNcz/CeBEQ=="

	c := NewMiniProgramClient(appId, appSecret, "")

	encryptedData := "98A19r2TH/F+biCYbx6YE9dnMjWVZfUgqPEWFenYX3jP8JpIKihNCyjE/Or0/pmYT+PCn6wCmV7s5LDwwQ92kcrMpuInOrmPWD36pI0mfywqh+53BcN4G+d30aG6ehCV3hPEqxE35ImpXxE5xuWqqsLX0YvgCgA5hLBGRWRPGiiXN4eSrLvCNI58BlC8VG16Iz6Z89NqAQ5WsQrWPjJqygZOsGnkvTFTwKzs6eM4jlmceWv4B37NxecCAwRkkHZRxla5mdluL5lwKovfH5feDQTg2Ui2/4Mc/raIh9tXV6lUqEkn8f4yJFhjXwJuhXdITxLcKnqP6O/n1DUOn1xNh5imZEiDBil14Zvy71pvb5JdRPhurtqOFM/5wCRc7BXwrrDo1iMUw3e7OOL9gFu7Fg=="
	iv := "+DTYOLIe/s1ObkVKZQolpw=="

	userInfo, err := c.DecryptUserInfo(sessionKey, encryptedData, iv)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%#v", userInfo)
}
