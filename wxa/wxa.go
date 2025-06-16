package wxa

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib/lib-wechat/wxa/api"
	"github.com/hdget/lib/lib-wechat/wxa/cache"
)

type Lib interface {
	Login(code string) (string, string, error) // 获取OpenId和UnionId
}

type wxmpImpl struct {
	api   api.Api
	cache cache.Cache
}

func New(appId, appSecret string, redisProvider intf.RedisProvider) Lib {
	return &wxmpImpl{
		api:   api.New(appId, appSecret),
		cache: cache.New(appId, redisProvider),
	}
}

func (impl wxmpImpl) Login(code string) (string, string, error) {
	result, err := impl.api.Code2Session(code)
	if err != nil {
		return "", "", err
	}

	// 保存到缓存中
	err = impl.cache.SetSessionKey(result.SessionKey)
	if err != nil {
		return "", "", err
	}

	return result.OpenId, result.UnionId, err
}
