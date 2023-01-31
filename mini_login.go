package goweixin

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"strings"
	"sync"
)

func init() {
	miner.SetLogLevel(miner.WARN)
}

type MiniProgramClient struct {
	AppId             string
	AppSecret         string
	AccessToken       string
	AccessTokenExpire int64
	accessTokenLock   sync.Mutex
	MiniProgramState  string
}

func NewMiniProgramClient(appId, appSecret string, programState string) *MiniProgramClient {
	m := new(MiniProgramClient)
	m.AppId = appId
	m.AppSecret = appSecret
	m.MiniProgramState = programState
	if programState == "" {
		m.MiniProgramState = MiniProgramStateFormal
	}
	return m
}

type MiniProgramBaseInfo struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

type MiniProgramUserInfo struct {
	NickName  string `json:"nickName"`
	OpenId    string `json:"openId"`
	AvatarUrl string `json:"avatarUrl"`
	UnionId   string `json:"unionId"`

	// 微信不再返回这部分信息了
	//Sex       int64                  `json:"sex"`
	//City      string                 `json:"city"`
	//Province  string                 `json:"province"`
	//Country   string                 `json:"country"`

	Watermark map[string]interface{} `json:"watermark"`

	SessionKey string `json:"session_key"`
}

func (c *MiniProgramClient) LoginGetBaseInfo(code string) (info *MiniProgramBaseInfo, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		c.AppId, c.AppSecret, code)

	api := miner.NewAPI()
	data, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return nil, err
	}

	miner.Logger.Infof("MiniProgramClient LoginGetSessionKey raw: %s", string(data))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(data, wErr)
	if err != nil {
		return
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	wToken := new(MiniProgramBaseInfo)
	err = json.Unmarshal(data, wToken)
	if err != nil {
		return
	}

	return wToken, nil

}

func (c *MiniProgramClient) LoginGetUserInfo(code, encryptedData, iv string) (info *MiniProgramUserInfo, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		c.AppId, c.AppSecret, code)

	api := miner.NewAPI()
	data, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return nil, err
	}

	miner.Logger.Infof("MiniProgramClient LoginGetUserInfo raw: %s", string(data))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(data, wErr)
	if err != nil {
		return
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	wToken := new(MiniProgramBaseInfo)
	err = json.Unmarshal(data, wToken)
	if err != nil {
		return
	}

	miner.Logger.Infof("MiniProgramClient LoginGetUserInfo token: %#v", wToken)

	sessionKey := wToken.SessionKey

	uInfo, err := c.DecryptUserInfo(sessionKey, encryptedData, iv)
	if err != nil {
		return
	}

	uInfo.OpenId = wToken.OpenId
	uInfo.UnionId = wToken.UnionId
	uInfo.SessionKey = sessionKey
	return uInfo, nil

}

// decryptWXOpenData https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/signature.html
func (c *MiniProgramClient) decryptWXOpenData(sessionKey, encryptData, iv string) ([]byte, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encryptData)
	if err != nil {
		return nil, err
	}

	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	dataBytes, err := AesDecrypt(encryptedData, sessionKeyBytes, ivBytes)
	if err != nil {
		return nil, err
	}

	return dataBytes, nil

}

func (c *MiniProgramClient) DecryptUserInfo(sessionKey, encryptData, iv string) (info *MiniProgramUserInfo, err error) {
	raw, err := c.decryptWXOpenData(sessionKey, encryptData, iv)
	if err != nil {
		return
	}

	// {\"nickName\":\"阿大\",\"gender\":0,\"language\":\"zh_CN\",\"city\":\"\",\"province\":\"\",\"country\":\"\",\"avatarUrl\":\"https://thirdwx.qlogo.cn/mmopen/vi_32/8JDU7pm9u0qUiaQr1ackQMA55RaE7Q3avOZ1YhG5kKsziaKM8YTDibtrH7rVicsRJu3YLlV0L3qG6YANa4dtF55zqA/132\",\"watermark\":{\"timestamp\":1669357949,\"appid\":\"wxd4e08529844845e7\"}}\x03\x03\x03
	miner.Logger.Infof("MiniProgramClient DecryptUserInfo raw %s", string(raw))

	temp := strings.Split(string(raw), "}")
	tempL := len(temp)
	if tempL < 2 {
		err = errors.New("not a json")
		return
	}

	temp2 := strings.Join(temp[:tempL-1], "}")
	raw = []byte(temp2 + "}")

	uInfo := new(MiniProgramUserInfo)
	err = json.Unmarshal(raw, uInfo)
	if err != nil {
		return
	}

	if uInfo.Watermark == nil {
		err = errors.New("watermark wrong nil")
		return
	}

	temp3, ok := uInfo.Watermark["appid"]
	if !ok {
		err = errors.New("watermark wrong app id not found")
		return
	}

	temp4 := fmt.Sprintf("%v", temp3)
	if temp4 != c.AppId {
		err = errors.New(fmt.Sprintf("watermark wrong app id not match, %s!=%s", temp4, c.AppId))
		return
	}

	return uInfo, nil
}

func AesDecrypt(encryptedData, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(origData, encryptedData)

	for i, ch := range origData {
		if ch == '\x0e' {
			origData[i] = ' '
		}
	}

	return origData, nil
}
