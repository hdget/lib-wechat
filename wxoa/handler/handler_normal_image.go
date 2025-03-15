package handler

import (
	"encoding/xml"
	"github.com/hdget/lib-wechat/wxoa"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type ImageMessage struct {
	*ReceivedMessage
	PicUrl    string
	MediaId   string
	Content   string
	MsgId     int64
	MsgDataId string
	Idx       string
}

type imageMessageHandler struct {
	callback wxoa.MessageCallback
	msg      *ImageMessage
}

func newTextMsgHandler(api wxoa.ApiWxoa, data []byte, callback) (wxoa.MessageHandler, error) {
	var m ImageMessage
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal image msg, data: %s", convert.BytesToString(data))
	}

	return &imageMessageHandler{instance: api, msg: &m}, nil
}

func (h imageMessageHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText(h.msg.Content)
}
