package provider

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

type WxaCode2SessionResult struct {
	*wxapi.Result
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

const (
	// 代商家管理小程序 /小程序登录 /小程序登录
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/login/thirdpartyCode2Session.html
	urlWxaCode2Session = "https://api.weixin.qq.com/sns/component/jscode2session?component_appid=%s&component_access_token=%s&appid=%s&js_code=%s&grant_type=authorization_code"
)

func (impl serviceProviderWxApiImpl) WxaCode2Session(componentAppId, componentAccessToken string, appId, code string) (*WxaCode2SessionResult, error) {
	url := fmt.Sprintf(urlWxaCode2Session,
		componentAppId,
		componentAccessToken,
		appId,
		code,
	)

	ret, err := wxapi.Post[WxaCode2SessionResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "wxa code to session")
	}

	if err = wxapi.CheckResult(ret.Result, url); err != nil {
		return nil, errors.Wrap(err, "wxa code to session")
	}

	return ret, nil
}
