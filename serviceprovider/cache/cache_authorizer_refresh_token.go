package cache

const (
	redisKeyAuthorizerRefreshToken = "authorizer_refresh_token"
)

func (c cacheImpl) GetAuthorizerRefreshToken(authorizerAppid string) (string, error) {
	return c.HGet(redisKeyAuthorizerRefreshToken, authorizerAppid)
}

func (c cacheImpl) SetAuthorizerRefreshToken(authorizerAppid string, refreshToken string) error {
	return c.HSet(redisKeyAuthorizerRefreshToken, authorizerAppid, refreshToken)
}
