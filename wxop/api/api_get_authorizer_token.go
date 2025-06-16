package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
)

type getAuthorizerAccessTokenRequest struct {
	ComponentAppid         string `json:"component_appid"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

type GetAuthorizerAccessTokenResult struct {
	api.Result
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	// 第三方平台调用凭证 /获取授权账号调用令牌 限制：2000次/天/授权方 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerAccessToken.html
	urlGetAuthorizerAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=%s"
)

func (impl apiImpl) GetAuthorizerAccessToken(componentAccessToken string, authorizerAppid, authorizerRefreshToken string) (*GetAuthorizerAccessTokenResult, error) {
	req := &getAuthorizerAccessTokenRequest{
		ComponentAppid:         impl.GetAppId(),
		AuthorizerAppid:        authorizerAppid,
		AuthorizerRefreshToken: authorizerRefreshToken,
	}

	url := fmt.Sprintf(urlGetAuthorizerAccessToken, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	var result GetAuthorizerAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return nil, err
	}

	if result.AuthorizerAccessToken == "" {
		return nil, fmt.Errorf("invalid authorizer access token result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return &result, nil
}
