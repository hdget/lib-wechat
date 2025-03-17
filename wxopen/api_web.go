package wxopen

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
)

// loginResult 类型
type loginResult struct {
	api.ApiResult
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	UnionId      string `json:"unionid"`
	Scope        string `json:"scope"`
}

const (
	// 参考: https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	urlOAuth2AccessToken = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
)

func (impl *wxopenImpl) Login(code string) (string, string, error) {
	url := fmt.Sprintf(urlOAuth2AccessToken, impl.AppId, impl.AppSecret, code)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return "", "", err
	}

	var result loginResult
	err = impl.ParseApiResult(resp.Body(), &result)
	if err != nil {
		return "", "", err
	}

	if result.AccessToken == "" {
		return "", "", fmt.Errorf("empty access token, url: %s, resp: %s", url, string(resp.Body()))
	}

	return result.OpenId, result.UnionId, nil
}
