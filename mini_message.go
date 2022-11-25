package goweixin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hunterhug/marmot/miner"
)

var (
	// MiniProgramStateDeveloper 开发版
	MiniProgramStateDeveloper = "developer"
	// MiniProgramStateTrial 体验版
	MiniProgramStateTrial = "trial"
	// MiniProgramStateFormal 正式版
	MiniProgramStateFormal = "formal"
)

type MiniProgramMessage struct {
	// 接收者（用户）的 openid
	ToUser string `json:"to_user"`
	// 所需下发的订阅模板id
	TemplateId string `json:"template_id"`
	// 点击模板卡片后的跳转页面，仅限本小程序内的页面。支持带参数,（示例index?foo=bar）。该字段不填则模板无跳转。
	Page             string `json:"page"`
	MiniProgramState string `json:"miniprogram_state"`
	Lang             string `json:"lang"`
	// 模板内容，格式形如 { "key1": { "value": any }, "key2": { "value": any } }
	Data map[string]interface{} `json:"data"`
}

type ErrorRsp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e *ErrorRsp) Error() string {
	return fmt.Sprintf("wxErr errcode:%d, errmsg:%v", e.ErrCode, e.ErrMsg)
}

// SendMessage https://developers.weixin.qq.com/miniprogram/dev/api/open-api/subscribe-message/wx.requestSubscribeMessage.html
func (c *MiniProgramClient) SendMessage(token string, openId string, templateId, page string, data map[string]string, state string) error {
	if token == "" || openId == "" || templateId == "" || state == "" {
		return errors.New("empty")
	}

	if state != MiniProgramStateDeveloper && state != MiniProgramStateFormal && state != MiniProgramStateTrial {
		return errors.New("state wrong")
	}

	// https://developers.weixin.qq.com/miniprogram/dev/api-backend/open-api/subscribe-message/subscribeMessage.send.html
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", token)
	m := new(MiniProgramMessage)
	m.ToUser = openId
	m.TemplateId = templateId
	m.Page = page
	m.Lang = "zh_CN"
	mm := map[string]interface{}{}

	for k, v := range data {
		mm[k] = map[string]string{"value": v}
	}

	m.Data = mm
	m.MiniProgramState = state

	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}

	worker := miner.NewAPI().Clone()
	body, err := worker.SetUrl(url).SetBData(raw).PostJSON()
	if err != nil {
		return err
	}

	miner.Logger.Infof("MiniProgramClient SendMessage result: %v", string(body))

	if worker.ResponseStatusCode != 200 {
		return errors.New(fmt.Sprintf("wx send message http status:%d", worker.ResponseStatusCode))
	}

	e := ErrorRsp{}
	err = json.Unmarshal(body, &e)
	if err != nil {
		return err
	}

	if e.ErrCode != 0 {
		return errors.New(e.ErrMsg)
	}
	return nil
}
