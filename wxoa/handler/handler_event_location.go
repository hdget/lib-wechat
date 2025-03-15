package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type locationEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newLocationEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &locationEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *locationEventMsgHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText("地理位置上报！")
}
