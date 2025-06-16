package wxopen

import "github.com/hdget/lib/lib-wechat/wxopen/api"

type Lib interface {
}

type wxopenImpl struct {
	api api.Api
}

func New(appId, appSecret string) Lib {
	return &wxopenImpl{
		api: api.New(appId, appSecret),
	}
}
