package wxop

import "encoding/xml"

type xmlComponentVerifyTicketEvent struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

const (
	redisKeyComponentVerifyTicket = "component_access_token"
)

func (impl wxopImpl) getComponentVerifyTicket() (string, error) {
	return impl.Cache.Get(redisKeyComponentVerifyTicket)
}

func (impl wxopImpl) apiProcessComponentVerifyTicket(data []byte) error {
	var event xmlComponentVerifyTicketEvent
	if err := xml.Unmarshal(data, &event); err != nil {
		return err
	}

	if err := impl.Cache.Set(redisKeyComponentVerifyTicket, event.ComponentVerifyTicket); err != nil {
		return err
	}

	return nil
}
