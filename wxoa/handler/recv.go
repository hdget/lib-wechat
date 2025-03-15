package handler

import (
	"encoding/xml"
	"github.com/hdget/lib-wechat/wxoa/reply"
	"github.com/pkg/errors"
	"time"
)

type ReceivedMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
}

func (m *ReceivedMessage) ReplyText(content string) ([]byte, error) {
	reply := reply.Text{
		XMLName:      xml.Name{},
		ToUserName:   m.FromUserName,
		FromUserName: m.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      content,
	}

	output, err := xml.MarshalIndent(reply, " ", " ")
	if err != nil {
		return nil, errors.Wrapf(err, "marshal text msg, reply: %v", reply)
	}

	return output, nil
}
