package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

const (
	urlGetAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

type GetAccessTokenResult struct {
	api.Result
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (impl apiImpl) GetAccessToken() (*GetAccessTokenResult, error) {
	url := fmt.Sprintf(urlGetAccessToken, impl.GetAppId(), impl.GetAppSecret())
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "get wxmp access token, appId: %s", impl.GetAppId())
	}

	var result GetAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if err = impl.CheckResult(result.Result, url, nil); err != nil {
		return nil, err
	}

	return &result, nil
}

//

//
//func (a *Api) GetAccessToken() (string, error) {
//	if a.Api == nil {
//		return "", errors.New("no redis provider")
//	}
//
//	// 尝试从缓存中获取access token
//	cachedAccessToken, _ := a.Api.Get("access_token")
//	if cachedAccessToken != "" {
//		return cachedAccessToken, nil
//	}
//
//	// 如果从缓存中获取不到，尝试请求access token
//	result, err := a.generateAccessToken()
//	if err != nil {
//		return "", err
//	}
//
//	err = a.Api.Set("access_token", result.AccessToken, result.ExpiresIn-1000)
//	if err != nil {
//		return "", err
//	}
//
//	return result.AccessToken, nil
//}
