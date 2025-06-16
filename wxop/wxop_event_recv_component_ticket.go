package wxop

import (
	"encoding/xml"
	"github.com/pkg/errors"
)

type xmlComponentVerifyTicketEvent struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

type componentVerifyTicketEventProcessor struct {
	*eventImpl
}

func (impl wxopImpl) newComponentVerifyTicketEventProcessor(e *eventImpl) EventProcessor {
	return &componentVerifyTicketEventProcessor{
		eventImpl: e,
	}
}

func (h componentVerifyTicketEventProcessor) Process(data []byte) (string, error) {
	var xmlEvent xmlComponentVerifyTicketEvent
	if err := xml.Unmarshal(data, &xmlEvent); err != nil {
		return "", err
	}

	return xmlEvent.ComponentVerifyTicket, nil

	// return cache.New(api.KindWxop, h.appId, h.redisProvider).Set(redisKeyComponentVerifyTicket, xmlEvent.ComponentVerifyTicket)
}

func (impl wxopImpl) updateComponentVerifyTicket(componentVerifyTicket string) error {
	if componentVerifyTicket == "" {
		return errors.New("empty component verify ticket")
	}
	return impl.cache.SetComponentVerifyTicket(componentVerifyTicket)
}
