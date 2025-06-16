package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hdget/lib/lib-wechat/api"
	"github.com/pkg/errors"
)

// GetSessionResult wechat miniprogram login session
type GetSessionResult struct {
	api.Result
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

const (
	// Code2Session 通过code换取openid和unionid
	// 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	urlLogin = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func (impl apiImpl) Code2Session(code string) (*GetSessionResult, error) {
	// 登录凭证校验
	url := fmt.Sprintf(urlLogin, impl.GetAppId(), impl.GetAppSecret(), code)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "wxmp login")
	}

	var result GetSessionResult
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return nil, err
	}

	if err = impl.CheckResult(result.Result, url, nil); err != nil {
		return nil, err
	}

	return &result, nil
}
