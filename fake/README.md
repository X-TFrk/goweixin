# 暂不使用

## 微信第三方登录

适用于网页端，移动端APP的微信登录。参考[官方文档](https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html) 。

需要客户端和服务端联调。

逻辑如下：

1.客户端先调用以下接口，微信用户允许授权第三方应用后，微信将会携带 `CODE` 并且回调服务端 `http://127.0.0.1:9999`：

https://open.weixin.qq.com/connect/qrconnect?appid=wx01fdsffsds&redirect_uri=http://127.0.0.1:9999&response_type=code&scope=snsapi_login,snsapi_userinfo&state=test#wechat_redirect

2.服务端收到回调，会连续调用以下链接获取到用户信息。

https://api.weixin.qq.com/sns/oauth2/access_token?appid=wx0189ce76eadccf91&secret=00cc512fc031fcdsfsdfba01c8a41f05b4b5&code=CODE&grant_type=authorization_code

https://api.weixin.qq.com/sns/userinfo?access_token=accessToken&openid=openid&lang=zh_CN

你只需使用该 `SDK` 实现登录即可：

```go
	appId := ""
	appSecret := ""
	code := "xxx" // 客户端传给你的，客户端可以是Web前端，IOS，Android
	info, err := Login(appId, appSecret, code)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(info)
```

