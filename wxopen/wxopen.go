package wxopen

import (
	"github.com/hdget/lib-wechat"
	"github.com/hdget/lib-wechat/api"
)

type wxopenImpl struct {
	*api.Api
}

var (
	_ ApiWxopen = (*wxopenImpl)(nil)
)

func New(appId, appSecret string) (ApiWxopen, error) {
	b, err := api.New(wechat.ApiKindWxopen, appId, appSecret)
	if err != nil {
		return nil, err
	}
	return &wxopenImpl{
		Api: b,
	}, nil
}
