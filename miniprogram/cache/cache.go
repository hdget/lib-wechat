package cache

import (
	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/pkg/cache"
)

type Cache interface {
	GetSessionKey() (string, error)
	SetSessionKey(sessionKey string) error
}

type cacheImpl struct {
	cache.Cache
}

func New(appId string, redis types.RedisProvider) Cache {
	return &cacheImpl{Cache: cache.New(cache.KindMiniProgram, appId, redis)}
}
