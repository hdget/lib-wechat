package api

type ApiKind string // 微信业务类型

const (
	ApiKindWxmp   ApiKind = "wxmp"   // 微信小程序
	ApiKindWxoa   ApiKind = "wxoa"   // 微信公众号
	ApiKindWxopen ApiKind = "wxopen" // 微信开放平台
	ApiKindWxop   ApiKind = "wxop"   // 微信第三方平台
)
