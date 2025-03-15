package wxmp

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat"
)

type wxmpImpl struct {
	*wechat.Api
}

var (
	_ ApiWxmp = (*wxmpImpl)(nil)
)

func New(appId, appSecret string, redisProvider intf.RedisProvider) (ApiWxmp, error) {
	a, err := wechat.New(wechat.ApiKindWxmp, appId, appSecret, redisProvider)
	if err != nil {
		return nil, err
	}

	return &wxmpImpl{
		Api: a,
	}, nil
}
