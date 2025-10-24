package miniprogram

import (
	"fmt"

	"github.com/hdget/lib-wechat/pkg/wxapi"
	"github.com/pkg/errors"
)

// GetSessionResult wechat miniprogram login session
type GetSessionResult struct {
	*wxapi.Result
	SessionKey string `json:"session_key"`
	OpenId     string `json:"openid"`
	UnionId    string `json:"unionid"`
}

const (
	// Code2Session 通过code换取openid和unionid
	// 参考：https://developers.weixin.qq.com/miniprogram/dev/OpenApiDoc/user-login/code2Session.html
	urlLogin = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

func (impl miniProgramWxApiImpl) Code2Session(code string) (*GetSessionResult, error) {
	// 登录凭证校验
	url := fmt.Sprintf(urlLogin, impl.GetAppId(), impl.GetAppSecret(), code)

	ret, err := wxapi.Get[GetSessionResult](url)
	if err != nil {
		return nil, errors.Wrap(err, "mini program code to session")
	}

	if err = wxapi.CheckResult(ret.Result, url); err != nil {
		return nil, errors.Wrap(err, "mini program code to session")
	}

	return ret, nil
}
