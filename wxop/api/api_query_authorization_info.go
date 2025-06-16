package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type queryAuthorizationInfoRequest struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

type queryAuthorizationInfoResult struct {
	api.Result
	AuthorizationInfo AuthorizationInfo `json:"authorization_info"`
}

type AuthorizationInfo struct {
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	// 第三方平台调用凭证/获取刷新令牌， 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getAuthorizerRefreshToken.html
	urlQueryAuthorizationInfo = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"
)

// QueryAuthorizationInfo 通过authCode获取AuthorizationInfo
func (impl apiImpl) QueryAuthorizationInfo(componentAccessToken, authCode string) (*AuthorizationInfo, error) {
	req := &queryAuthorizationInfoRequest{
		ComponentAppid:    impl.GetAppId(),
		AuthorizationCode: authCode,
	}

	url := fmt.Sprintf(urlQueryAuthorizationInfo, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	var result queryAuthorizationInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal authorization info, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return nil, err
	}

	if result.AuthorizationInfo.AuthorizerAccessToken == "" {
		return nil, fmt.Errorf("invalid authorizer access token result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return &result.AuthorizationInfo, nil
}
