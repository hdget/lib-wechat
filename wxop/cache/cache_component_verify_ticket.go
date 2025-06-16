package cache

const (
	redisKeyComponentVerifyTicket = "component_verify_ticket"
)

func (c cacheImpl) GetComponentVerifyTicket() (string, error) {
	return c.Get(redisKeyComponentVerifyTicket)
}

func (c cacheImpl) SetComponentVerifyTicket(componentVerifyTicket string) error {
	return c.Set(redisKeyComponentVerifyTicket, componentVerifyTicket)
}
