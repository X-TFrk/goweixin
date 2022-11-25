package goweixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/miner"
)

type MiniProgramPhoneNumber struct {
	Info *MiniProgramPhoneInfo `json:"phone_info"`
}

type MiniProgramPhoneInfo struct {
	PhoneNumber     string                 `json:"phoneNumber"`
	PurePhoneNumber string                 `json:"purePhoneNumber"`
	CountryCode     string                 `json:"countryCode"`
	Watermark       map[string]interface{} `json:"watermark"`
}

// GetPhoneNumber https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/phonenumber/phonenumber.getPhoneNumber.html
func (c *MiniProgramClient) GetPhoneNumber(token string, code string) (*MiniProgramPhoneInfo, error) {
	if token == "" || code == "" {
		return nil, errors.New("empty")
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", token)

	worker := miner.NewAPI().Clone()
	body, err := worker.SetUrl(url).SetBData([]byte(fmt.Sprintf(`{ "code": "%s" }`, code))).PostJSON()
	if err != nil {
		return nil, err
	}

	miner.Logger.Infof("MiniProgramClient GetPhoneNumber raw: %s", string(body))

	wErr := new(ErrorRsp)
	err = json.Unmarshal(body, wErr)
	if err != nil {
		return nil, err
	}

	if wErr.ErrCode != 0 {
		return nil, wErr
	}

	// {"errcode":0,"errmsg":"ok","phone_info":{"phoneNumber":"1322214333","purePhoneNumber":"13221433","countryCode":"86","watermark":{"timestamp":1669367712,"appid":"wx8c6b032dbc5cd756"}}}
	phoneNumber := new(MiniProgramPhoneNumber)
	err = json.Unmarshal(body, phoneNumber)
	if err != nil {
		return nil, err
	}

	if phoneNumber.Info == nil || phoneNumber.Info.Watermark == nil {
		err = errors.New("watermark wrong nil")
		return nil, err
	}

	temp3, ok := phoneNumber.Info.Watermark["appid"]
	if !ok {
		err = errors.New("watermark wrong app id not found")
		return nil, err
	}

	temp4 := fmt.Sprintf("%v", temp3)
	if temp4 != c.AppId {
		err = errors.New(fmt.Sprintf("watermark wrong app id not match, %s!=%s", temp4, c.AppId))
		return nil, err
	}

	return phoneNumber.Info, nil
}
