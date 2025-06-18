package cache

const (
	sessionKeyExpires  = 3600 // session key过期时间3600秒
	redisKeySessionKey = "session_key"
)

func (c *cacheImpl) GetSessionKey() (string, error) {
	return c.Get(redisKeySessionKey)
}

func (c *cacheImpl) SetSessionKey(sessionKey string) error {
	return c.Set(redisKeySessionKey, sessionKey, sessionKeyExpires)
}
