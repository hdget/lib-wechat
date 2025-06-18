package web

import (
	"github.com/hdget/lib-wechat/op/web/api"
)

type Lib interface {
	Login(code string) (string, string, error) // 获取OpenId和UnionId
}

type webAppImpl struct {
	api api.Api
}

func New(appId, appSecret string) Lib {
	return &webAppImpl{
		api: api.New(appId, appSecret),
	}
}

func (impl webAppImpl) Login(code string) (string, string, error) {
	return impl.api.Login(code)
}
