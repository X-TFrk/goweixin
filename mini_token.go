package goweixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/miner"
)

type MiniProgramAccessToken struct {
	AccessToken string `json:"access_token"`
}

// AuthGetAccessToken https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/access-token/auth.getAccessToken.html
func (c *MiniProgramClient) AuthGetAccessToken() (token string, err error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", c.AppId, c.AppSecret)
	api := miner.NewAPI()
	raw, err := api.Clone().SetUrl(url).Get()
	if err != nil {
		return "", err
	}

	miner.Logger.Infof("MiniProgramClient AuthGetAccessToken token: %v", string(raw))

	t := new(MiniProgramAccessToken)
	err = json.Unmarshal(raw, t)
	if err != nil {
		return "", err
	}

	if t.AccessToken == "" {
		return "", errors.New("empty")
	}

	return t.AccessToken, nil
}
