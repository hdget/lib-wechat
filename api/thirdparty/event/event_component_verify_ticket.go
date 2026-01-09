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

type componentVerifyTicketEventHandler struct {
	appId         string
	redisProvider types.RedisProvider
}

func newComponentVerifyTicketEventHandler() Handler {
	return &componentVerifyTicketEventHandler{}
}

func (h componentVerifyTicketEventHandler) Handle(data []byte) (string, error) {
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
