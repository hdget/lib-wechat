package wxop

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib-wechat/api"
	"github.com/pkg/errors"
)

type getPreAuthCodeRequest struct {
	ComponentAppid string `json:"component_appid"`
}

type getPreAuthCodeResult struct {
	api.ApiResult
	PreAuthCode string `json:"pre_auth_code"`
	ExpiresIn   int    `json:"expires_in"`
}

const (
	urlGetPreAuthCode = "https://api.weixin.qq.com/cgi-bin/component/api_create_preauthcode?component_access_token=%s"
)

func (impl wxopImpl) getPreAuthCode() (string, error) {
	componentAccessToken, err := impl.getComponentAccessToken()
	if err != nil {
		return "", errors.Wrap(err, "get component verify ticket")
	}

	req := &getPreAuthCodeRequest{
		ComponentAppid: impl.AppId,
	}

	url := fmt.Sprintf(urlGetPreAuthCode, componentAccessToken)
	resp, err := resty.New().R().SetBody(req).Post(url)
	if err != nil {
		return "", err
	}

	var result getPreAuthCodeResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("%s, url: %s", result.ErrMsg, url)
	}

	if result.PreAuthCode == "" {
		return "", fmt.Errorf("invalid pre auth code result, url: %s, resp: %s", url, string(resp.Body()))
	}

	return result.PreAuthCode, nil
}
