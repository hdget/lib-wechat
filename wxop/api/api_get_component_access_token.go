package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type getComponentAccessTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type GetComponentAccessTokenResult struct {
	api.Result
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证 /获取令牌 限制：2000次/天
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
	urlGetComponentAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
)

func (impl apiImpl) GetComponentAccessToken(componentVerifyTicket string) (*GetComponentAccessTokenResult, error) {
	req := &getComponentAccessTokenRequest{
		ComponentAppid:        impl.GetAppId(),
		ComponentAppsecret:    impl.GetAppSecret(),
		ComponentVerifyTicket: componentVerifyTicket,
	}

	resp, err := resty.New().R().SetBody(req).Post(urlGetComponentAccessToken)
	if err != nil {
		return nil, err
	}

	var result GetComponentAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal component access token result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, urlGetComponentAccessToken, req); err != nil {
		return nil, err
	}

	if result.ComponentAccessToken == "" {
		return nil, fmt.Errorf("invalid component access token result, url: %s, resp: %s", urlGetComponentAccessToken, string(resp.Body()))
	}

	return &result, nil
}
