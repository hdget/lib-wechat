package event

import (
	"encoding/xml"

	"github.com/hdget/common/types"
	"github.com/hdget/lib-wechat/api/thirdparty/cache"
	"github.com/pkg/errors"
)

type xmlComponentVerifyTicketEvent struct {
	ComponentVerifyTicket string `xml:"ComponentVerifyTicket"`
}

type componentVerifyTicketEventProcessor struct {
	appId         string
	redisProvider types.RedisProvider
}

func newComponentVerifyTicketEventProcessor() PreProcessor {
	return &componentVerifyTicketEventProcessor{}
}

func (h componentVerifyTicketEventProcessor) Process(data []byte) (string, error) {
	var e xmlComponentVerifyTicketEvent
	if err := xml.Unmarshal(data, &e); err != nil {
		return "", err
	}

	if e.ComponentVerifyTicket == "" {
		return "", errors.New("empty component verify ticket")
	}

	_ = cache.ComponentVerifyTicket(h.appId, h.redisProvider).Set(e.ComponentVerifyTicket)

	return e.ComponentVerifyTicket, nil
}
