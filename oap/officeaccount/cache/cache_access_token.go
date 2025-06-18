package cache

const (
	redisKeyAccessToken = "access_token"
)

func (c cacheImpl) GetAccessToken() (string, error) {
	return c.Get(redisKeyAccessToken)
}

func (c cacheImpl) SetAccessToken(accessToken string, expiresIn int) error {
	return c.Set(redisKeyAccessToken, accessToken, expiresIn)
}
