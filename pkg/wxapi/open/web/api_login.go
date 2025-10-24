package web

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

// loginResult 类型
type loginResult struct {
	*wxapi.Result
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
func (impl openWebWxApiImpl) Login(code string) (string, string, error) {
	url := fmt.Sprintf(urlOAuth2AccessToken, impl.GetAppId(), impl.GetAppSecret(), code)

	ret, err := wxapi.Get[loginResult](url)
	if err != nil {
		return "", "", errors.Wrap(err, "open platform web login")
	}

	if err = wxapi.CheckResult(ret.Result, url); err != nil {
		return "", "", errors.Wrap(err, "open platform web login")
	}

	return ret.OpenId, ret.UnionId, nil
}
