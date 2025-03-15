package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type unsubScanEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newUnsubScanEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &unsubScanEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *newUnsubScanEventMsgHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText("未关注用户扫码！")
}
