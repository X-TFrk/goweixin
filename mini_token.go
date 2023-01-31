package goweixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/miner"
	"time"
)

type MiniProgramAccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// AuthGetAccessToken https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func (c *MiniProgramClient) AuthGetAccessToken() (token string, err error) {
	c.accessTokenLock.Lock()
	defer c.accessTokenLock.Unlock()

	if c.AccessToken != "" && c.AccessTokenExpire <= time.Now().Unix() {
		return c.AccessToken, nil
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", c.AppId, c.AppSecret)
	api := miner.NewAPI()
	raw, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return "", err
	}

	miner.Logger.Infof("MiniProgramClient AuthGetAccessToken token: %v", string(raw))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(raw, wErr)
	if err != nil {
		return "", err
	}

	if wErr.ErrCode != 0 {
		return "", wErr
	}

	t := new(MiniProgramAccessToken)
	err = json.Unmarshal(raw, t)
	if err != nil {
		return "", err
	}

	if t.AccessToken == "" {
		return "", errors.New("token empty")
	}

	c.AccessToken = t.AccessToken
	c.AccessTokenExpire = time.Now().Unix() + t.ExpiresIn - 5
	return t.AccessToken, nil
}
