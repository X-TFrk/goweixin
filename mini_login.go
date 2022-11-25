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
)

type miniAccessToken struct {
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

type MiniUserInfo struct {
	NickName  string                 `json:"nickName"`
	OpenId    string                 `json:"openId"`
	AvatarUrl string                 `json:"avatarUrl"`
	UnionId   string                 `json:"unionId"`
	Sex       int64                  `json:"sex"`
	City      string                 `json:"city"`
	Province  string                 `json:"province"`
	Country   string                 `json:"country"`
	Watermark map[string]interface{} `json:"watermark"`
}

func MiniLogin(appId, appSecret, code, encryptedData, iv string) (info *MiniUserInfo, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		appId, appSecret, code)

	api := miner.NewAPI()
	data, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return nil, err
	}

	wErr := new(ErrorRsp)
	err = json.Unmarshal(data, wErr)
	if err != nil {
		return
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	wToken := new(miniAccessToken)
	err = json.Unmarshal(data, wToken)
	if err != nil {
		return
	}

	miner.Logger.Infof("wx MiniLogin token: %#v", wToken)

	sessionKey := wToken.SessionKey

	raw, err := DecryptWXOpenData(sessionKey, encryptedData, iv)
	if err != nil {
		return
	}

	miner.Logger.Infof("wx MiniLogin get userInfo: %#v", string(data))

	temp := strings.Split(string(raw), "}")
	tempL := len(temp)
	if tempL < 2 {
		err = errors.New("not a json")
		return
	}

	temp2 := strings.Join(temp[:tempL-1], "}")
	raw = []byte(temp2 + "}")

	uInfo := new(MiniUserInfo)
	err = json.Unmarshal(raw, uInfo)
	if err != nil {
		return
	}

	uInfo.OpenId = wToken.OpenId
	uInfo.UnionId = wToken.UnionId

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
	if temp4 != appId {
		err = errors.New(fmt.Sprintf("watermark wrong app id not match, %s!=%s", temp4, appId))
		return
	}

	return uInfo, nil

}

// DecryptWXOpenData https://developers.weixin.qq.com/miniprogram/dev/framework/open-ability/signature.html
func DecryptWXOpenData(sessionKey, encryptData, iv string) ([]byte, error) {
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
