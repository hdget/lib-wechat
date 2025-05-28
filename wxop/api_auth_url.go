package wxop

import (
	"fmt"
	"github.com/spf13/cast"
	"net/url"
	"strings"
)

const (
	urlPCAuth = "https://mp.weixin.qq.com/cgi-bin/componentloginpage?component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%s"
	urlH5Auth = "https://open.weixin.qq.com/wxaopen/safe/bindcomponent?action=bindcomponent&no_scan=1&component_appid=%s&pre_auth_code=%s&redirect_uri=%s&auth_type=%s#wechat_redirect"
)

func (impl wxopImpl) GetAuthUrl(client, redirectUrl, authCode string) (string, error) {
	preAuthCode, err := impl.getPreAuthCode()
	if err != nil {
		return "", err
	}

	// 校验authCode
	switch cast.ToInt(authCode) {
	case 1, 2, 3, 4, 5, 6:
	default:
		return "", fmt.Errorf("invalid auth code: %s", authCode)
	}

	// 校验redirectUrl, https://xxx
	if !impl.isValidRedirectURL(redirectUrl) {
		return "", fmt.Errorf("invalid redirect url, redirectUrl: %s", redirectUrl)
	}

	switch strings.ToLower(client) {
	case "pc":
		return fmt.Sprintf(urlPCAuth, impl.AppId, preAuthCode, redirectUrl, authCode), nil
	case "h5":
		return fmt.Sprintf(urlH5Auth, impl.AppId, preAuthCode, redirectUrl, authCode), nil
	default:
		return "", fmt.Errorf("unsupported client, client: %s", client)
	}
}

func (impl wxopImpl) isValidRedirectURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	// 检查是否有Scheme
	if u.Scheme == "" {
		return false
	}

	// 检查Scheme是否合法（https）
	if u.Scheme != "https" {
		return false
	}

	// 检查Host是否为空
	if u.Host == "" {
		return false
	}

	// 可选：检查是否有路径或查询参数
	// if u.Path == "" && u.RawQuery == "" {
	// 	return false
	// }

	return true
}
