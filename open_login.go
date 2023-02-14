package goweixin

import (
	"encoding/json"
	"fmt"
	"github.com/hunterhug/marmot/miner"
)

type OpenClient struct {
	AppId     string
	AppSecret string
}

type OpenBaseInfo struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
}

type OpenUserInfo struct {
	NickName string `json:"nickname"`
	OpenId   string `json:"openid"`
	Img      string `json:"headimgurl"`
	UnionId  string `json:"unionid"`
	Sex      int64  `json:"sex"`
	City     string `json:"city"`
	Province string `json:"province"`
	Country  string `json:"country"`
}

func NewOpenClient(appId, appSecret string) *OpenClient {
	m := new(OpenClient)
	m.AppId = appId
	m.AppSecret = appSecret
	return m
}

func (c *OpenClient) LoginGetBaseInfo(code string) (info *OpenBaseInfo, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		c.AppId, c.AppSecret, code)

	api := miner.NewAPI()
	data, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return nil, err
	}

	miner.Logger.Infof("OpenClient LoginGetBaseInfo raw: %s", string(data))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(data, wErr)
	if err != nil {
		return
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	wToken := new(OpenBaseInfo)
	err = json.Unmarshal(data, wToken)
	if err != nil {
		return
	}

	return wToken, nil
}

func (c *OpenClient) LoginGetUserInfo(code string) (info *OpenUserInfo, err error) {
	wToken, err := c.LoginGetBaseInfo(code)
	if err != nil {
		return nil, err
	}

	accessToken := wToken.AccessToken
	openId := wToken.OpenId

	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, openId)
	api := miner.NewAPI()
	data, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return nil, err
	}

	miner.Logger.Infof("OpenClient LoginGetUserInfo raw: %s", string(data))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(data, wErr)
	if err != nil {
		return
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	uInfo := new(OpenUserInfo)
	err = json.Unmarshal(data, uInfo)
	if err != nil {
		return
	}

	return uInfo, nil
}
