package wxoa

import "github.com/hdget/lib-wechat/wxoa/handler"

// HandleRecv 处理接收到的消息
func (impl *wxoaImpl) HandleRecv(data []byte) ([]byte, error) {
	h, err := handler.New(impl, data)
	if err != nil {
		return nil, err
	}
	return h.Handle()
}
