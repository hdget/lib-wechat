package provider

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type getAuthorizerAccessTokenRequest struct {
	ComponentAppid         string `json:"component_appid"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

type GetAuthorizerAccessTokenResult struct {
	*wxapi.Result
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	// 第三方平台调用凭证 /获取授权账号调用令牌 限制：2000次/天/授权方 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerAccessToken.html
	urlGetAuthorizerAccessToken = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=%s"
)

func (impl serviceProviderWxApiImpl) GetAuthorizerAccessToken(componentAccessToken string, authorizerAppid, authorizerRefreshToken string) (*GetAuthorizerAccessTokenResult, error) {
	req := &getAuthorizerAccessTokenRequest{
		ComponentAppid:         impl.GetAppId(),
		AuthorizerAppid:        authorizerAppid,
		AuthorizerRefreshToken: authorizerRefreshToken,
	}

	url := fmt.Sprintf(urlGetAuthorizerAccessToken, componentAccessToken)

	ret, err := wxapi.Post[GetAuthorizerAccessTokenResult](url, req)
	if err != nil {
		return nil, errors.Wrap(err, "get authorizer access token")
	}

	if err = wxapi.CheckResult(ret.Result, url, req); err != nil {
		return nil, errors.Wrap(err, "get authorizer access token")
	}

	if ret.AuthorizerAccessToken == "" {
		return nil, errors.New("empty authorizer access token")
	}

	return ret, nil
}
