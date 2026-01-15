package event

import (
	"encoding/xml"
	"strings"

	"github.com/hdget/lib-wechat/pkg/crypt"
	"github.com/pkg/errors"
)

type AuthEventKind string

const (
	AuthEventKindComponentVerifyTicket AuthEventKind = "component_verify_ticket" // 校验组件校验凭证
	AuthEventKindAuthorized            AuthEventKind = "authorized"              // 授权
	AuthEventKindUpdateAuthorized      AuthEventKind = "updateauthorized"        // 更新授权
	AuthEventKindUnauthorized          AuthEventKind = "unauthorized"            // 取消授权
)

type Event interface {
	Handle() error
}

// Message 接收到的消息
type Message struct {
	Signature string
	Timestamp string
	Nonce     string
	Body      string
}

type AuthEventHandler func(string) error // 处理事件

// PreProcessor 预处理器，对消息进行预处理
type PreProcessor interface {
	Process(data []byte) (string, error)
}

type authEventImpl struct {
	kind AuthEventKind
	data []byte
}

var (
	// 授权事件预处理器
	_authEventPreProcessors = map[AuthEventKind]PreProcessor{
		AuthEventKindComponentVerifyTicket: newComponentVerifyTicketEventProcessor(),
		AuthEventKindAuthorized:            newAuthorizedEventProcessor(),
		AuthEventKindUnauthorized:          newUnauthorizedEventProcessor(),
		AuthEventKindUpdateAuthorized:      newAuthorizedEventProcessor(),
	}
	_authEventHandlers = map[AuthEventKind]AuthEventHandler{}
)

type xmlAuthEvent struct {
	InfoType string `xml:"InfoType"`
}

// NewAuthEvent 创建授权事件
func NewAuthEvent(appId, token, encodingAESKey string, message *Message) (Event, error) {
	msgCrypt, err := crypt.NewBizMsgCrypt(appId, token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	data, err := msgCrypt.Decrypt(message.Signature, message.Timestamp, message.Nonce, message.Body)
	if err != nil {
		return nil, err
	}

	var evt xmlAuthEvent
	if err = xml.Unmarshal(data, &evt); err != nil {
		return nil, err
	}

	return &authEventImpl{
		kind: AuthEventKind(strings.ToLower(evt.InfoType)),
		data: data,
	}, nil
}

// RegisterAuthEventHandler 注册授权事件处理Handler
func RegisterAuthEventHandler(kind AuthEventKind, handler AuthEventHandler) {
	_authEventHandlers[kind] = handler
}

func (impl authEventImpl) Handle() error {
	var result string
	var err error
	if preProcessor, exists := _authEventPreProcessors[impl.kind]; exists {
		result, err = preProcessor.Process(impl.data)
		if err != nil {
			return errors.Wrapf(err, "pre process event, kind: %s, data: %s", impl.kind, string(impl.data))
		}
	}

	handler, exists := _authEventHandlers[impl.kind]
	if !exists {
		return errors.Wrapf(err, "handler not exists, kind: %s, handlers: %v", impl.kind, _authEventHandlers)
	}

	if err = handler(result); err != nil {
		return errors.Wrapf(err, "handle event, kind: %s, result: %s", impl.kind, result)
	}

	return nil
}
