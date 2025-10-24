package provider

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type createPreAuthCodeRequest struct {
	ComponentAppid string `json:"component_appid"`
}

type createPreAuthCodeResult struct {
	*wxapi.Result
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证/获取预授权码 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
	urlCreatePreAuthCode = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=%s"
)

func (impl serviceProviderWxApiImpl) CreatePreAuthCode(componentAccessToken string) (string, error) {
	req := &createPreAuthCodeRequest{
		ComponentAppid: impl.GetAppId(),
	}

	url := fmt.Sprintf(urlCreatePreAuthCode, componentAccessToken)

	ret, err := wxapi.Post[createPreAuthCodeResult](url, req)
	if err != nil {
		return "", errors.Wrap(err, "get authorizer option")
	}

	if err = wxapi.CheckResult(ret.Result, url, req); err != nil {
		return "", errors.Wrap(err, "get authorizer option")
	}

	if ret.PreAuthCode == "" {
		return "", errors.New("empty pre auth code")
	}

	return ret.PreAuthCode, nil
}
