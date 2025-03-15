package wxopen

import (
	"github.com/hdget/lib-wechat"
)

type wxopenImpl struct {
	*wechat.Api
}

var (
	_ ApiWxopen = (*wxopenImpl)(nil)
)

func New(appId, appSecret string) (ApiWxopen, error) {
	b, err := wechat.New(wechat.ApiKindWxopen, appId, appSecret)
	if err != nil {
		return nil, err
	}
	return &wxopenImpl{
		Api: b,
	}, nil
}
