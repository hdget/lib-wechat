package wxop

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

type eventImpl struct {
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

func (impl wxopImpl) HandleEvent(event *Event, handlers map[EventKind]EventHandler) error {
	crypt, err := NewWXBizMsgCrypt(impl.api.GetAppId(), event.Token, event.EncodingAESKey)
	if err != nil {
		return err
	}

	data, err := crypt.DecryptMsg(event.Signature, event.Timestamp, event.Nonce, event.Body)
	if err != nil {
		return err
	}

	var e eventImpl
	if err = xml.Unmarshal(data, &e); err != nil {
		return err
	}

	infoType := EventKind(strings.ToLower(e.InfoType))

	var processor EventProcessor
	handler := handlers[infoType]
	switch infoType {
	case EventKindComponentVerifyTicket:
		processor = impl.newComponentVerifyTicketEventProcessor(&e)
		if handler == nil {
			handler = impl.updateComponentVerifyTicket
		}
	case EventKindAuthorized, EventKindUpdateAuthorized: // 授权成功, 换取authorizer info
		processor = impl.newAuthorizedEventProcessor(&e)
	case EventKindUnAuthorized: // 取消授权
		processor = impl.newUnauthorizedEventProcessor(&e)
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
