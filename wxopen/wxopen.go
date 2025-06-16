package wxopen

import "github.com/hdget/lib/lib-wechat/wxopen/api"

type Lib interface {
	Login(code string) (string, string, error) // 获取OpenId和UnionId
}

type wxopenImpl struct {
	api api.Api
}

func New(appId, appSecret string) Lib {
	return &wxopenImpl{
		api: api.New(appId, appSecret),
	}
}

func (impl wxopenImpl) Login(code string) (string, string, error) {
	return impl.api.Login(code)
}
