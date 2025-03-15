package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type subScanEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newSubScanEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &subScanEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *subScanEventMsgHandler) Handle() ([]byte, error) {
	return h.msg.ReplyText("已关注用户扫码！")
}
