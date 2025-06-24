package serviceprovider

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"strings"
)

type EventKind string

const (
	EventKindComponentVerifyTicket EventKind = "component_verify_ticket"
	EventKindAuthorized            EventKind = "authorized"
	EventKindUpdateAuthorized      EventKind = "updateauthorized"
	EventKindUnAuthorized          EventKind = "unauthorized"
)

type EventProcessor interface {
	Process(data []byte) (string, error) // 预处理
}

// EventHandler 授权回调函数
type EventHandler func(input string) error

type xmlEvent struct {
	InfoType string `xml:"InfoType"`
}

type Event struct {
	Token          string
	EncodingAESKey string
	Signature      string
	Timestamp      string
	Nonce          string
	Body           string
}

func (impl serviceProviderImpl) HandleEvent(event *Event, handlers map[EventKind]EventHandler) error {
	crypt, err := NewWXBizMsgCrypt(impl.api.GetAppId(), event.Token, event.EncodingAESKey)
	if err != nil {
		return err
	}

	data, err := crypt.DecryptMsg(event.Signature, event.Timestamp, event.Nonce, event.Body)
	if err != nil {
		return err
	}

	var evt xmlEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return err
	}

	infoType := EventKind(strings.ToLower(evt.InfoType))

	var processor EventProcessor
	handler := handlers[infoType]
	switch infoType {
	case EventKindComponentVerifyTicket:
		processor = impl.newComponentVerifyTicketEventProcessor(&evt)
		if handler == nil {
			handler = impl.updateComponentVerifyTicket
		}
	case EventKindAuthorized, EventKindUpdateAuthorized: // 授权成功, 换取authorizer info
		processor = impl.newAuthorizedEventProcessor(&evt)
	case EventKindUnAuthorized: // 取消授权
		processor = impl.newUnauthorizedEventProcessor(&evt)
	default:
		return errors.New("unsupported event processor")
	}

	result, err := processor.Process(data)
	if err != nil {
		return err
	}

	if handler != nil {
		return handler(result)
	}

	return nil
}
