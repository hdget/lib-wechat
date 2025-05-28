package wxopen

import (
	"github.com/hdget/lib-wechat/api"
)

type ApiWxopen interface {
	/*
	 * 网站应用快速扫码登录
	 * 参考：https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Wechat_Login.html
	 */
	Login(code string) (string, string, error) // 网站应用扫码快速登录
}

type wxopenImpl struct {
	*api.Api
}

var (
	_ ApiWxopen = (*wxopenImpl)(nil)
)

func New(appId, appSecret string) (ApiWxopen, error) {
	b, err := api.New(api.ApiKindWxopen, appId, appSecret)
	if err != nil {
		return nil, err
	}
	return &wxopenImpl{
		Api: b,
	}, nil
}
