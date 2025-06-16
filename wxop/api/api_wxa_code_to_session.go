package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type WxaCode2SessionResult struct {
	api.Result
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

const (
	// 代商家管理小程序 /小程序登录 /小程序登录
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/miniprogram-management/login/thirdpartyCode2Session.html
	urlWxaCode2Session = "https://api.weixin.qq.com/sns/component/jscode2session?component_appid=%s&component_access_token=%s&appid=%s&js_code=%s&grant_type=authorization_code"
)

func (impl apiImpl) WxaCode2Session(componentAppId, componentAccessToken string, appId, code string) (*WxaCode2SessionResult, error) {
	url := fmt.Sprintf(urlWxaCode2Session,
		componentAppId,
		componentAccessToken,
		appId,
		code,
	)

	resp, err := resty.New().R().Post(url)
	if err != nil {
		return nil, err
	}

	var result WxaCode2SessionResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal wxa code2session result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, url); err != nil {
		return nil, err
	}

	return &result, nil
}
