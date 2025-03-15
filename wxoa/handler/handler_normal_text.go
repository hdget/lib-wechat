package handler

import (
	"encoding/xml"
	"github.com/hdget/lib-wechat/wxoa"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type TextMessage struct {
	*ReceivedMessage
	Content   string
	MsgId     int64
	MsgDataId string
	Idx       string
}

type textMessageHandler struct {
	callback wxoa.MessageCallback
	msg      *TextMessage
}

func newTextMsgHandler(api wxoa.ApiWxoa, data []byte) (wxoa.MessageHandler, error) {
	var m TextMessage
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal event msg, data: %s", convert.BytesToString(data))
	}

	return &textMessageHandler{instance: api, msg: &m}, nil
}

func (h textMessageHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText(h.msg.Content)
}
