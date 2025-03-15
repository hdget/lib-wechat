package wechat

type ApiKind int // 微信业务类型

const (
	ApiKindUnknown ApiKind = iota
	ApiKindWxmp            // 微信小程序
	ApiKindWxoa            // 微信公众号
	ApiKindWxopen          // 微信开放平台
)

// ApiError 微信的错误响应
type ApiError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
