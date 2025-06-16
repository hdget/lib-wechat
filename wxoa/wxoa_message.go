package wxoa

import (
	"github.com/hdget/lib/lib-wechat/wxoa/message"
	"sync"
)

var (
	locker              sync.Mutex
	_msgKind2msgHandler = map[message.MessageKind]message.MessageHandler{}
)

// HandleMessage 处理消息
func (impl *wxoaImpl) HandleMessage(data []byte) ([]byte, error) {
	m, err := message.New(data)
	if err != nil {
		return nil, err
	}

	if h, exists := _msgKind2msgHandler[m.GetKind()]; exists {
		return h(m)
	}
	return m.Reply()
}

func RegisterMessageHandler(msgKind message.MessageKind, handler message.MessageHandler) error {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler[msgKind] = handler
	return nil
}
