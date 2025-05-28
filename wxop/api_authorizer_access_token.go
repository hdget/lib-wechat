package wxop

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type getAuthorizerAccessTokenRequest struct {
	ComponentAppid         string `json:"component_appid"`
	AuthorizerAppid        string `json:"authorizer_appid"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

type getAuthorizerAccessTokenResult struct {
	api.ApiResult
	AuthorizerAccessToken  string `json:"authorizer_access_token"`
	ExpiresIn              int    `json:"expires_in"`
	AuthorizerRefreshToken string `json:"authorizer_refresh_token"`
}

const (
	urlGetAuthorizerAccessToken    = "https://api.weixin.qq.com/cgi-bin/component/api_authorizer_token?component_access_token=%s"
	redisKeyAuthorizerAccessToken  = "authorizer_access_token:%s"
	redisKeyAuthorizerRefreshToken = "authorizer_access_token"
)

func (impl wxopImpl) GetAuthorizerAccessToken(authorizerAppid string, fnLoadAuthRefreshToken LoadAuthorizerRefreshToken) (string, error) {
	key := fmt.Sprintf(redisKeyAuthorizerAccessToken, authorizerAppid)
	authAccessToken, err := impl.Cache.Get(key)
	if err != nil {
		authRefreshToken, err := impl.getAuthorizerRefreshToken(authorizerAppid, fnLoadAuthRefreshToken)
		if err != nil {
			return "", err
		}

		result, err := impl.apiGetAuthorizerAccessToken(authorizerAppid, authRefreshToken)
		if err != nil {
			return "", err
		}

		// cache auth_access_token和auth_refresh_token
		err = impl.Cache.Set(key, result.AuthorizerAccessToken, result.ExpiresIn-600)
		if err != nil {
			return "", errors.Wrap(err, "cache authorizer access token")
		}

		err = impl.Cache.HSet(redisKeyAuthorizerRefreshToken, authorizerAppid, result.AuthorizerRefreshToken)
		if err != nil {
			return "", errors.Wrap(err, "cache authorizer refresh token")
		}

		return result.AuthorizerAccessToken, nil
	}
	return authAccessToken, nil
}

func (impl wxopImpl) apiGetAuthorizerAccessToken(authorizerAppid, authorizerRefreshToken string) (*getAuthorizerAccessTokenResult, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return nil, errors.Wrap(err, "get component verify ticket")
	}

	req := &getAuthorizerAccessTokenRequest{
		ComponentAppid:         impl.AppId,
		AuthorizerAppid:        authorizerAppid,
		AuthorizerRefreshToken: authorizerRefreshToken,
	}

	url := fmt.Sprintf(urlGetAuthorizerAccessToken, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return nil, err
	}

	var result getAuthorizerAccessTokenResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("%s, url: %s", result.ErrMsg, url)
	}

	if result.AuthorizerAccessToken == "" {
		return nil, fmt.Errorf("invalid authorizer access token result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return &result, nil
}

// getAuthorizerRefreshToken 获取保存的authorizerRefreshToken, 先从缓存中找，找不到从持久层中找
func (impl wxopImpl) getAuthorizerRefreshToken(authorizerAppid string, fnLoadAuthRefreshToken LoadAuthorizerRefreshToken) (string, error) {
	authRefreshToken, err := impl.Cache.HGet(redisKeyAuthorizerRefreshToken, authorizerAppid)
	if err != nil {
		authRefreshToken, err = fnLoadAuthRefreshToken(authorizerAppid)
		if err != nil {
			return "", errors.Wrapf(err, "load authorizer refresh token, authorizerAppid: %s", authorizerAppid)
		}

		err = impl.Cache.HSet(redisKeyAuthorizerRefreshToken, authorizerAppid, authRefreshToken)
		if err != nil {
			return "", errors.Wrap(err, "cache authorizer refresh token")
		}

		return authRefreshToken, nil
	}
	return authRefreshToken, nil
}
