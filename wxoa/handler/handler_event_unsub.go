package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type unsubEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newUnsubEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &unsubEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *unsubEventMsgHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText("取消关注！")
}
