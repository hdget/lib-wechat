package authevent

import (
	"encoding/xml"

	"github.com/hdget/lib-wechat/pkg/crypt"
)

type AppEventKind string

type AppEventHandler func() error

type AppEvent interface {
	RegisterHandler(kind AuthEventKind, handler AppEventHandler)
	Handle() error
}

type appEventImpl struct {
	kind     AuthEventKind
	data     []byte
	handlers map[AuthEventKind]AppEventHandler
}

type xmlAppEvent struct {
	ToUserName string `xml:"ToUserName"`
	Encrypt    string `xml:"Encrypt"`
}

func NewAppEvent(appId, token, encodingAESKey string, message *Message) (AppEvent, error) {
	msgCrypt, err := crypt.NewBizMsgCrypt(appId, token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	data, err := msgCrypt.Decrypt(message.Signature, message.Timestamp, message.Nonce, message.Body)
	if err != nil {
		return nil, err
	}

	var evt xmlAppEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return nil, err
	}

	return &appEventImpl{
		data:     data,
		handlers: make(map[AuthEventKind]AppEventHandler),
	}, nil
}

// RegisterHandler 注册代运营APP事件处理Handler
func (impl appEventImpl) RegisterHandler(kind AuthEventKind, handler AppEventHandler) {
	impl.handlers[kind] = handler
}

func (impl appEventImpl) Handle() error {
	if handler, ok := impl.handlers[impl.kind]; ok {
		if err := handler(); err != nil {
			return err
		}
	}

	return nil
}
