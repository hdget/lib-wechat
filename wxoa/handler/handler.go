package handler

import (
	"encoding/xml"
	"fmt"
	"github.com/hdget/lib-wechat/wxoa"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type baseHandler struct {
	callback wxoa.MessageCallback
	msg      any
}

func (b *baseHandler) Handle() ([]byte, error) {
	return b.callback.Execute(b.msg)
}

func (b *baseHandler) getCallback(h wxoa.MessageHandler, callbacks ...wxoa.MessageCallback) wxoa.MessageCallback {
	if len(callbacks) == 0 {
		return h.GetDefaultCallback()
	}
	return callbacks[0]
}

type newHandlerFunction func([]byte, wxoa.MessageCallback) (wxoa.MessageHandler, error)

var (
	msgType2fn = map[string]newHandlerFunction{
		"text":  newTextMsgHandler,
		"image": newTextMessageHandler,
		"event": newEventHandler,
	}
)

func New(api wxoa.ApiWxoa, data []byte, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	var msg ReceivedMessage
	err := xml.Unmarshal(data, &msg)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal msg, data: %s", convert.BytesToString(data))
	}

	if fn, exists := msgType2fn[msg.MsgType]; exists {
		return fn(data, cb)
	}

	return nil, fmt.Errorf("unsupported msg type, msgType: %s", msg.MsgType)
}
