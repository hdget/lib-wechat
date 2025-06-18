package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
)

// GetJsSdkTicketResult 类型
type GetJsSdkTicketResult struct {
	api.Result
	Value     string `json:"ticket,omitempty"`
	ExpiresIn int    `json:"expires_in,omitempty"`
}

const (
	// 参考：Js SDK签名：https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html
	urlGetJsSdkTicket = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// ApiGetJsSdkTicket jssdk获取凭证
func (impl apiImpl) GetJsSdkTicket(accessToken string) (*GetJsSdkTicketResult, error) {
	url := fmt.Sprintf(urlGetJsSdkTicket, accessToken)
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}

	var result GetJsSdkTicketResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if err = impl.CheckResult(result.Result, url, nil); err != nil {
		return nil, err
	}

	return &result, nil
}
