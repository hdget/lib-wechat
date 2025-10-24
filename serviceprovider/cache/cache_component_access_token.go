package cache

const (
	redisKeyComponentAccessToken = "component_access_token"
)

func (c cacheImpl) GetComponentAccessToken() (string, error) {
	return c.Get(redisKeyComponentAccessToken)
}

func (c cacheImpl) SetComponentAccessToken(accessToken string, expiresIn int) error {
	return c.Set(redisKeyComponentAccessToken, accessToken, expiresIn)
}
