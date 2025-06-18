package officeaccount

import (
	message2 "github.com/hdget/lib-wechat/oap/officeaccount/message"
	"sync"
)

var (
	locker              sync.Mutex
	_msgKind2msgHandler = map[message2.MessageKind]message2.MessageHandler{}
)

// HandleMessage 处理消息
func (impl *wxoaImpl) HandleMessage(data []byte) ([]byte, error) {
	m, err := message2.New(data)
	if err != nil {
		return nil, err
	}

	if h, exists := _msgKind2msgHandler[m.GetKind()]; exists {
		return h(m)
	}
	return m.Reply()
}

func RegisterMessageHandler(msgKind message2.MessageKind, handler message2.MessageHandler) error {
	locker.Lock()
	defer locker.Unlock()
	_msgKind2msgHandler[msgKind] = handler
	return nil
}
