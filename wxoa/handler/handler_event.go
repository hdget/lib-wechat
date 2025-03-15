package handler

import (
	"encoding/xml"
	"fmt"
	"github.com/hdget/lib-wechat/wxoa"
	"github.com/hdget/utils/convert"
	"github.com/pkg/errors"
)

type EventMessage struct {
	*ReceivedMessage
	Event     string
	EventKey  string
	Ticket    string
	Latitude  float64
	Longitude float64
	Precision float64
}

func newEventHandler(data []byte, cb wxoa.MessageCallback) (wxoa.MessageHandler, error) {
	var eventMessage EventMessage
	err := xml.Unmarshal(data, &eventMessage)
	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal event msg, data: %s", convert.BytesToString(data))
	}

	switch eventMessage.Event {
	case "subscribe":
		if eventMessage.Ticket != "" { // 未关注用户扫码
			return newUnsubScanEventMsgHandler(&eventMessage, cb)
		} else { // 关注公众号
			return newSubEventMsgHandler(&eventMessage, cb)
		}
	case "unsubscribe": // 取消关注公众号
		return newUnsubEventMsgHandler(&eventMessage, cb)
	case "SCAN": // 已关注用户扫码
		return newSubScanEventMsgHandler(&eventMessage, cb)
	case "LOCATION":
		return newLocationEventMsgHandler(&eventMessage, cb)
	case "CLICK":
		return newClickEventMsgHandler(&eventMessage, cb)
	case "VIEW":
		return newViewEventMsgHandler(&eventMessage, cb)
	}

	return nil, fmt.Errorf("unsupported event message, event: %s", eventMessage.Event)
}
