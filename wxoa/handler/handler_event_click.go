package handler

import (
	"github.com/hdget/lib-wechat/wxoa"
)

type clickEventMsgHandler struct {
	*baseHandler
}

func newClickEventMsgHandler(eventMsg *EventMessage, callbacks ...wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	h := &clickEventMsgHandler{
		baseHandler: &baseHandler{
			msg: eventMsg,
		},
	}

	h.baseHandler.callback = h.getCallback(h, callbacks...)
	return h, nil
}

func (h *clickEventMsgHandler) GetDefaultCallback() ([]byte, error) {
	return h.msg.ReplyText("点击事件")
}
