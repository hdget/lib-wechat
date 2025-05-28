package wxop

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type queryAuthorizationInfoRequest struct {
	ComponentAppid    string `json:"component_appid"`
	AuthorizationCode string `json:"authorization_code"`
}

type QueryAuthorizationInfoResult struct {
	api.ApiResult
	AuthorizationInfo AuthorizationInfo `json:"authorization_info"`
}

type AuthorizationInfo struct {
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	urlQueryAuthorizationInfo = "https://api.weixin.qq.com/cgi-bin/component/api_query_auth?component_access_token=%s"
)

func (impl wxopImpl) apiQueryAuthorizationInfo(authorizationCode string) (*AuthorizationInfo, error) {
	componentVerifyTicket, err := impl.GetComponentVerifyTicket()
	if err != nil {
		return nil, errors.Wrap(err, "get component verify ticket")
	}

	req := &queryAuthorizationInfoRequest{
		ComponentAppid:    impl.AppId,
		AuthorizationCode: authorizationCode,
	}

	url := fmt.Sprintf(urlQueryAuthorizationInfo, componentVerifyTicket)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	var result QueryAuthorizationInfoResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, url)
	}

	if result.AuthorizationInfo.AuthorizerAccessToken == "" {
		return nil, fmt.Errorf("invalid authorizer access token result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return &result.AuthorizationInfo, nil
}
