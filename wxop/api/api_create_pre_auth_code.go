package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type createPreAuthCodeRequest struct {
	ComponentAppid string `json:"component_appid"`
}

type createPreAuthCodeResult struct {
	api.Result
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	// 第三方平台调用凭证/获取预授权码 限制：2000次/天/平台
	// https://developers.weixin.qq.com/doc/oplatform/openApi/OpenApiDoc/ticket-token/getPreAuthCode.html
	urlGetPreAuthCode = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=%s"
)

func (impl apiImpl) CreatePreAuthCode(componentAccessToken string) (string, error) {
	req := &createPreAuthCodeRequest{
		ComponentAppid: impl.GetAppId(),
	}

	url := fmt.Sprintf(urlGetPreAuthCode, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return "", err
	}

	var result createPreAuthCodeResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", errors.Wrapf(err, "unmarshal pre auth code result, data: %s", convert.BytesToString(resp.Body()))
	}

	if err = impl.CheckResult(result.Result, url, req); err != nil {
		return "", err
	}

	if result.PreAuthCode == "" {
		return "", fmt.Errorf("invalid pre auth code result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return result.PreAuthCode, nil
}
