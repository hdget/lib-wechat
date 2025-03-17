package wxoa

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat/api"
)

type wxoaImpl struct {
	*api.Api
}

var (
	_ ApiWxoa = (*wxoaImpl)(nil)
)

func New(appId, appSecret string, redisProvider intf.RedisProvider) (ApiWxoa, error) {
	a, err := api.New(api.ApiKindWxoa, appId, appSecret, redisProvider)
	if err != nil {
		return nil, err
	}

	return &wxoaImpl{
		Api: a,
	}, nil
}
