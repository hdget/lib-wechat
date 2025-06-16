package cache

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib/lib-wechat/cache"
)

type Cache interface {
	GetJsSdkTicket() (string, error)
	SetJsSdkTicket(ticket string, expiresIn int) error
	SetAccessToken(accessToken string, expiresIn int) error
	GetAccessToken() (string, error)
}

type cacheImpl struct {
	cache.Cache
}

func New(appId string, redis intf.RedisProvider) Cache {
	return &cacheImpl{Cache: cache.New(cache.KindWxoa, appId, redis)}
}
