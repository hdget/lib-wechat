package wxmp

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat"
	"github.com/hdget/lib-wechat/api"
)

type wxmpImpl struct {
	*api.Api
}

var (
	_ ApiWxmp = (*wxmpImpl)(nil)
)

func New(appId, appSecret string, redisProvider intf.RedisProvider) (ApiWxmp, error) {
	a, err := api.New(wechat.ApiKindWxmp, appId, appSecret, redisProvider)
	if err != nil {
		return nil, err
	}

	return &wxmpImpl{
		Api: a,
	}, nil
}
