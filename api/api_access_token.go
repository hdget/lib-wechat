package api

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type AccessTokenResult struct {
	Result
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	urlGetWechatAccessToken = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

func (a *Api) GetAccessToken() (string, error) {
	if a.Cache == nil {
		return "", errors.New("no redis provider")
	}

	// 尝试从缓存中获取access token
	cachedAccessToken, _ := a.Cache.Get("access_token")
	if cachedAccessToken != "" {
		return cachedAccessToken, nil
	}

	// 如果从缓存中获取不到，尝试请求access token
	result, err := a.generateAccessToken()
	if err != nil {
		return "", err
	}

	err = a.Cache.Set("access_token", result.AccessToken, result.ExpiresIn-1000)
	if err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

func (a *Api) generateAccessToken() (*AccessTokenResult, error) {
	wxAccessTokenURL := fmt.Sprintf(urlGetWechatAccessToken, a.AppId, a.AppSecret)

	resp, err := resty.New().R().Get(wxAccessTokenURL)
	if err != nil {
		return nil, errors.Wrapf(err, "get access token, appId: %s", a.AppId)
	}

	var result AccessTokenResult
	err = a.ParseApiResult(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.AccessToken == "" {
		return nil, fmt.Errorf("empty access token, url: %s, resp: %s", wxAccessTokenURL, convert.BytesToString(resp.Body()))
	}

	return &result, nil
}
