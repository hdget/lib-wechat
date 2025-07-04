package cache

import (
	"fmt"
)

const (
	redisKeyAuthorizerAccessToken = "authorizer_access_token:%s"
)

func (c cacheImpl) GetAuthorizerAccessToken(authorizerAppid string) (string, error) {
	return c.Get(fmt.Sprintf(redisKeyAuthorizerAccessToken, authorizerAppid))
}

func (c cacheImpl) SetAuthorizerAccessToken(authorizerAppid string, accessToken string, expiresIn int) error {
	return c.Set(fmt.Sprintf(redisKeyAuthorizerAccessToken, authorizerAppid), accessToken, expiresIn)
}
