# 微信服务端API接口Golang SDK

# 如何使用

```go
go get -u -v github.com/hunterhug/goweixin
```

## 一. 小程序开发

小程序基础库版本：[2.26.1](https://developers.weixin.qq.com/miniprogram/dev/framework/release)

### A. [小程序登录](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/login.html)

小程序登录区别于网页登录，需要客户端和服务端联调，获取密钥对请登陆 [微信公众平台](https://mp.weixin.qq.com)。

逻辑如下：

1. 客户端先调用 `wx.login()` 获取临时登录凭证 `code` 并且 [获取用户信息](https://developers.weixin.qq.com/miniprogram/dev/api/open-api/user-info/wx.getUserProfile.html) 获取 `encryptedData` 和 `iv` 并回传到开发者服务器。
2. 服务端使用该 `code` 调用 [`auth.code2Session`](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/login/auth.code2Session.html) 获取解密密钥，然后解密用户信息。

你只需使用该 `SDK` 实现登录并获取用户信息即可：

```
func TestMiniProgramClient_LoginGetUserInfo(t *testing.T) {
	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

	c := NewMiniProgramClient(appId, appSecret)

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
```

### B. 小程序发送 [消息订阅](https://developers.weixin.qq.com/miniprogram/dev/api/open-api/subscribe-message/wx.requestSubscribeMessage.html)

完全在服务端执行，不需要客户端参与。

1. 先获取全局 [access_token](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html) 。
2. 然后发送[订阅消息](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html) 。

```go
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

	openId := "on9VO5YXH_gMLxRKMql98IUjtzkI"
	templateId := "IgOxNz7ydQn9UsssswIwiggdd5TroAbRLXcosY"
	page := "index?foo=bar"
	data := map[string]string{"thing1": "这是一个内容", "thing7": "这个也是内容", "thing3": "这个也是内容啊"}

	err = c.SendMessage(token, openId, templateId, page, data)
	if err != nil {
		fmt.Println("send err:", err.Error())
		return
	}
}
```

### C. [获取手机号](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/getPhoneNumber.html)

需要客户端和服务端联调。

逻辑如下：

1. 客户端调用 [`<button open-type="getPhoneNumber" bindgetphonenumber="getPhoneNumber"></button>`](https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/getPhoneNumber.html) 获取到 `code` 传给服务端。
2. 服务端先获取全局 [access_token](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html) 。
3. 服务端再配合 `code` 调用 [phonenumber.getPhoneNumber](https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/phonenumber/phonenumber.getPhoneNumber.html) 获取手机号。

```go
func TestMiniProgramClient_GetPhoneNumber(t *testing.T) {
	appId := "wxd4e08529844845e7"
	appSecret := "e6782244f7a7e994d20721f004e3e9ae"

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
```