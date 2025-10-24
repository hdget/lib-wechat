package web

import (
	"github.com/hdget/lib-wechat/pkg/wxapi"
)

type WxAPI interface {
	Login(code string) (string, string, error)
}

type openWebWxApiImpl struct {
	wxapi.API
}

func New(appId, appSecret string) WxAPI {
	return &openWebWxApiImpl{
		API: wxapi.New(appId, appSecret),
	}
}
