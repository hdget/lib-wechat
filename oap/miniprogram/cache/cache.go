package cache

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib-wechat/cache"
)

type Cache interface {
	GetSessionKey() (string, error)
	SetSessionKey(sessionKey string) error
}

type cacheImpl struct {
	cache.Cache
}

func New(appId string, redis intf.RedisProvider) Cache {
	return &cacheImpl{Cache: cache.New(cache.KindWxa, appId, redis)}
}
