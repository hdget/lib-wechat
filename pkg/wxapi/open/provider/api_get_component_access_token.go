package provider

import (
	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type getComponentAccessTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type GetComponentAccessTokenResult struct {
	*wxapi.Result
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证 /获取令牌 限制：2000次/天
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getComponentAccessToken.html
	urlGetComponentAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
)

func (impl serviceProviderWxApiImpl) GetComponentAccessToken(componentVerifyTicket string) (*GetComponentAccessTokenResult, error) {
	req := &getComponentAccessTokenRequest{
		ComponentAppid:        impl.GetAppId(),
		ComponentAppsecret:    impl.GetAppSecret(),
		ComponentVerifyTicket: componentVerifyTicket,
	}

	ret, err := wxapi.Post[GetComponentAccessTokenResult](urlGetComponentAccessToken, req)
	if err != nil {
		return nil, errors.Wrap(err, "get component access token")
	}

	if err = wxapi.CheckResult(ret.Result, urlGetComponentAccessToken, req); err != nil {
		return nil, errors.Wrap(err, "get component access token")
	}

	if ret.ComponentAccessToken == "" {
		return nil, errors.New("empty component access token")
	}

	return ret, nil
}
