package event

import (
	"encoding/xml"
	"strings"

	"github.com/hdget/lib-wechat/pkg/crypt"
	"github.com/pkg/errors"
)

type Kind string

const (
	EventKindComponentVerifyTicket Kind = "component_verify_ticket"
	EventKindAuthorized            Kind = "authorized"
	EventKindUpdateAuthorized      Kind = "updateauthorized"
	EventKindUnAuthorized          Kind = "unauthorized"
)

type Event interface {
	Handle() (string, error)
}

type Message struct {
	Token          string
	EncodingAESKey string
	Signature      string
	Timestamp      string
	Nonce          string
	Body           string
}

type Handler interface {
	Handle(data []byte) (string, error) // 处理事件
}

type eventImpl struct {
	kind Kind
	data []byte
}

var (
	_handlers = map[Kind]Handler{
		EventKindComponentVerifyTicket: newComponentVerifyTicketEventHandler(),
		EventKindAuthorized:            newAuthorizedEventHandler(),
		EventKindUnAuthorized:          newUnauthorizedEventHandler(),
		EventKindUpdateAuthorized:      newAuthorizedEventHandler(),
	}
)

type xmlEvent struct {
	InfoType string `xml:"InfoType"`
}

func New(appId string, message *Message) (Event, error) {
	msgCrypt, err := crypt.NewBizMsgCrypt(appId, message.Token, message.EncodingAESKey)
	if err != nil {
		return nil, err
	}

	data, err := msgCrypt.Decrypt(message.Signature, message.Timestamp, message.Nonce, message.Body)
	if err != nil {
		return nil, err
	}

	var evt xmlEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return nil, err
	}

	return &eventImpl{
		kind: Kind(strings.ToLower(evt.InfoType)),
		data: data,
	}, nil
}

func RegisterHandler(kind Kind, handler Handler) {
	_handlers[kind] = handler
}

func (impl eventImpl) Handle() (string, error) {
	handler, exists := _handlers[impl.kind]
	if !exists {
		return "", errors.New("unsupported event processor")
	}

	return handler.Handle(impl.data)
}
