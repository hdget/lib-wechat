package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type viewEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newViewEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &viewEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *viewEventMsgHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText("跳转链接事件！")
}
