package handler

import (
	"fmt"
	"github.com/hdget/lib-wechat/wxoa"
	"github.com/pkg/errors"
)

type subEventMsgHandler struct {
	callback wxoa.MessageCallback
	msg      *EventMessage
}

func newSubEventMsgHandler(eventMsg *EventMessage, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	return &subEventMsgHandler{instance: instance, msg: eventMsg}, nil
}

func (h *subEventMsgHandler) Handle() ([]byte, error) {
	userInfo, err := h.instance.GetUserInfo(h.msg.FromUserName)
	if err != nil {
		return nil, errors.Wrapf(err, "wxoa get user info, fromUserName: %s", h.msg.FromUserName)
	}

	if userInfo.UnionId == "" {
		return nil, fmt.Errorf("union id not found, fromUserName: %s", h.msg.FromUserName)
	}

	err = h.instance.FollowWxoa(userInfo.UnionId, h.msg.FromUserName)
	if err != nil {
		return nil, errors.Wrapf(err, "follow wxoa, fromUserName: %s, unionid: %s", h.msg.FromUserName, userInfo.UnionId)
	}

	return h.msg.ReplyText("欢迎关注！")
}
