package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
)

// loginResult 类型
type loginResult struct {
	api.Result
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

// Login 网站应用快速扫码登录
func (impl apiImpl) Login(code string) (string, string, error) {
	url := fmt.Sprintf(urlOAuth2AccessToken, impl.GetAppId(), impl.GetAppSecret(), code)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return "", "", err
	}

	var result loginResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", "", err
	}

	if err = impl.CheckResult(result.Result, url, nil); err != nil {
		return "", "", err
	}

	return result.OpenId, result.UnionId, nil
}
