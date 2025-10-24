package cache

import (
	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/pkg/cache"
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

func New(appId string, redis types.RedisProvider) Cache {
	return &cacheImpl{Cache: cache.New(cache.KindOfficeAccount, appId, redis)}
}
