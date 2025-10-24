package cache

const (
	redisKeyJsApiTicket = "js_api_ticket"
)

func (c cacheImpl) GetJsSdkTicket() (string, error) {
	return c.Get(redisKeyJsApiTicket)
}

func (c cacheImpl) SetJsSdkTicket(ticket string, expiresIn int) error {
	return c.Set(redisKeyJsApiTicket, ticket, expiresIn)
}
