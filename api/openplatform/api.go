package openplatform

import (
	"github.com/hdget/lib-wechat/api"
	"github.com/hdget/lib-wechat/api/openplatform/wx"
)

type API interface {
	api.Api
	WebAppLogin(code string) (string, string, error) // 网站应用快速扫码登录
}

type openPlatformApiImpl struct {
	api.Api
	wx.WxApi
}

func New(appId, appSecret string) API {
	return &openPlatformApiImpl{
		Api:   api.New(appId, appSecret),
		WxApi: wx.New(appId, appSecret),
	}
}

func (impl openPlatformApiImpl) WebAppLogin(code string) (string, string, error) {
	return impl.WxApi.WebAppLogin(code)
}
