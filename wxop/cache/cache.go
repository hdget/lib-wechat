package cache

import (
	"github.com/hdget/common/intf"
	"github.com/hdget/lib/lib-wechat/cache"
)

type Cache interface {
	GetComponentVerifyTicket() (string, error)
	SetComponentVerifyTicket(componentVerifyTicket string) error
	GetComponentAccessToken() (string, error)
	SetComponentAccessToken(accessToken string, expiresIn int) error
	GetAuthorizerAccessToken(authorizerAppid string) (string, error)
	SetAuthorizerAccessToken(authorizerAppid string, accessToken string, expiresIn int) error
	GetAuthorizerRefreshToken(authorizerAppid string) (string, error)
	SetAuthorizerRefreshToken(authorizerAppid string, refreshToken string) error
}

type cacheImpl struct {
	cache.Cache
}

func New(appId string, redis intf.RedisProvider) Cache {
	return &cacheImpl{Cache: cache.New(cache.KindWxop, appId, redis)}
}
