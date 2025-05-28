package wxop

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type getComponentAccessTokenRequest struct {
	ComponentAppid        string `json:"component_appid"`
	ComponentAppsecret    string `json:"component_appsecret"`
	ComponentVerifyTicket string `json:"component_verify_ticket"`
}

type getComponentAccessTokenResult struct {
	api.ApiResult
	ComponentAccessToken string `json:"component_access_token"`
	ExpiresIn            int    `json:"expires_in"`
}

const (
	urlGetComponentAccessToken   = "https://api.weixin.qq.com/cgi-bin/component/api_component_token"
	redisKeyComponentAccessToken = "wxop:component_access_token"
)

func (impl wxopImpl) GetComponentAccessToken() (string, error) {
	componentAccessToken, err := impl.Cache.Get(redisKeyComponentAccessToken)
	if err != nil { // 缓存取不到则通过API接口获取并缓存起来
		result, err := impl.apiGetComponentAccessToken()
		if err != nil {
			return "", errors.Wrap(err, "retrieve component access token")
		}

		// 过期前十分钟过期
		err = impl.Cache.Set(redisKeyComponentAccessToken, result.ComponentAccessToken, result.ExpiresIn-600)
		if err != nil {
			return "", err
		}

		return result.ComponentAccessToken, nil
	}

	return componentAccessToken, nil
}

func (impl wxopImpl) apiGetComponentAccessToken() (*getComponentAccessTokenResult, error) {
	componentVerifyTicket, err := impl.GetComponentVerifyTicket()
	if err != nil {
		return nil, errors.Wrap(err, "get component verify ticket")
	}

	req := &getComponentAccessTokenRequest{
		ComponentAppid:        impl.AppId,
		ComponentAppsecret:    impl.AppSecret,
		ComponentVerifyTicket: componentVerifyTicket,
	}

	resp, err := resty.New().R().SetBody(req).Post(urlGetComponentAccessToken)
	if err != nil {
		return nil, err
	}

	var result getComponentAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, urlGetComponentAccessToken)
	}

	if result.ComponentAccessToken == "" {
		return nil, fmt.Errorf("invalid component access token result, url: %s, resp: %s", urlGetComponentAccessToken, string(resp.Body()))
	}

	return &result, nil
}
